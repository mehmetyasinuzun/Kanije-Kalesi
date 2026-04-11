// Package telegram implements the Telegram Bot API client.
// All strings are UTF-8. Message text uses Telegram's MarkdownV2 format.
// No third-party Telegram SDK is used — only stdlib net/http.
package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	apiBase    = "https://api.telegram.org/bot"
	maxRetries = 3
)

// Client is a low-level Telegram Bot API HTTP client.
// All methods automatically retry on transient errors with exponential back-off.
type Client struct {
	token      string
	httpClient *http.Client
	log        *slog.Logger
}

// NewClient creates a Telegram API client with sensible defaults.
func NewClient(token string, timeoutSec int, log *slog.Logger) *Client {
	if timeoutSec <= 0 {
		timeoutSec = 15
	}
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSec) * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: 90 * time.Second,
			},
		},
		log: log,
	}
}

// apiURL builds the full API endpoint URL.
func (c *Client) apiURL(method string) string {
	return apiBase + c.token + "/" + method
}

// ---- API response types ----

type apiResponse struct {
	OK          bool            `json:"ok"`
	Result      json.RawMessage `json:"result"`
	Description string          `json:"description"`
	ErrorCode   int             `json:"error_code"`
	Parameters  *retryParams    `json:"parameters"`
}

type retryParams struct {
	RetryAfter int `json:"retry_after"`
}

// Update is a single incoming update from Telegram.
type Update struct {
	UpdateID      int64          `json:"update_id"`
	Message       *Message       `json:"message"`
	CallbackQuery *CallbackQuery `json:"callback_query"`
}

// Message is a Telegram message.
type Message struct {
	MessageID int64  `json:"message_id"`
	From      *User  `json:"from"`
	Chat      *Chat  `json:"chat"`
	Text      string `json:"text"`
	Date      int64  `json:"date"`
}

// CallbackQuery represents a button press from an inline keyboard.
type CallbackQuery struct {
	ID      string   `json:"id"`
	From    *User    `json:"from"`
	Message *Message `json:"message"`
	Data    string   `json:"data"` // Callback data from the button
}

// User represents a Telegram user or bot.
type User struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

// Chat represents a Telegram chat.
type Chat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

// InlineKeyboardMarkup is a grid of inline keyboard buttons.
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// InlineKeyboardButton is a single button in an inline keyboard.
type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data,omitempty"`
	URL          string `json:"url,omitempty"`
}

// ---- Core API methods ----

// GetMe validates the bot token and returns bot information.
func (c *Client) GetMe(ctx context.Context) (*User, error) {
	var user User
	if err := c.call(ctx, "getMe", nil, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// SendMessage sends a UTF-8 text message.
// parseMode: "MarkdownV2", "HTML", or "" (plain text).
func (c *Client) SendMessage(ctx context.Context, chatID int64, text, parseMode string) (*Message, error) {
	params := map[string]any{
		"chat_id": chatID,
		"text":    text,
	}
	if parseMode != "" {
		params["parse_mode"] = parseMode
	}

	var msg Message
	if err := c.call(ctx, "sendMessage", params, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// SendMessageWithKeyboard sends a message with an inline keyboard.
func (c *Client) SendMessageWithKeyboard(ctx context.Context, chatID int64, text string, kb InlineKeyboardMarkup) (*Message, error) {
	kbJSON, _ := json.Marshal(kb)
	params := map[string]any{
		"chat_id":      chatID,
		"text":         text,
		"parse_mode":   "HTML",
		"reply_markup": string(kbJSON),
	}

	var msg Message
	if err := c.call(ctx, "sendMessage", params, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

// EditMessageReplyMarkup updates the inline keyboard of an existing message.
func (c *Client) EditMessageReplyMarkup(ctx context.Context, chatID, messageID int64, kb InlineKeyboardMarkup) error {
	kbJSON, _ := json.Marshal(kb)
	params := map[string]any{
		"chat_id":      chatID,
		"message_id":   messageID,
		"reply_markup": string(kbJSON),
	}
	return c.call(ctx, "editMessageReplyMarkup", params, nil)
}

// EditMessageText updates the text of an existing message.
func (c *Client) EditMessageText(ctx context.Context, chatID, messageID int64, text string, kb *InlineKeyboardMarkup) error {
	params := map[string]any{
		"chat_id":    chatID,
		"message_id": messageID,
		"text":       text,
		"parse_mode": "HTML",
	}
	if kb != nil {
		kbJSON, _ := json.Marshal(kb)
		params["reply_markup"] = string(kbJSON)
	}
	return c.call(ctx, "editMessageText", params, nil)
}

// SendPhoto sends a JPEG photo with an optional caption.
func (c *Client) SendPhoto(ctx context.Context, chatID int64, data []byte, caption string) error {
	return c.sendMedia(ctx, "sendPhoto", "photo", "photo.jpg", chatID, data, caption)
}

// AnswerCallbackQuery acknowledges a callback query (clears the button's loading state).
func (c *Client) AnswerCallbackQuery(ctx context.Context, callbackID, text string) error {
	params := map[string]any{
		"callback_query_id": callbackID,
	}
	if text != "" {
		params["text"] = text
	}
	return c.call(ctx, "answerCallbackQuery", params, nil)
}

// GetUpdates performs long-polling and returns new updates.
// timeout is the long-poll timeout in seconds (0 = short poll).
func (c *Client) GetUpdates(ctx context.Context, offset, timeout int64) ([]Update, error) {
	params := map[string]any{
		"offset":          offset,
		"timeout":         timeout,
		"allowed_updates": []string{"message", "callback_query"},
	}

	var updates []Update
	if err := c.call(ctx, "getUpdates", params, &updates); err != nil {
		return nil, err
	}
	return updates, nil
}

// ---- HTTP layer ----

// call makes an API request with automatic retry on rate limiting and transient errors.
func (c *Client) call(ctx context.Context, method string, params map[string]any, dest any) error {
	var lastErr error
	delay := 2 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
				delay *= 2
			}
		}

		resp, err := c.doRequest(ctx, method, params)
		if err != nil {
			lastErr = err
			continue
		}

		if !resp.OK {
			if resp.ErrorCode == 429 && resp.Parameters != nil {
				// Rate limited — respect Retry-After
				ra := time.Duration(resp.Parameters.RetryAfter) * time.Second
				c.log.Warn("Telegram hız sınırı", "bekleme", ra)
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(ra):
					attempt-- // Don't count this as a retry
					continue
				}
			}
			return fmt.Errorf("telegram API hatası %d: %s", resp.ErrorCode, resp.Description)
		}

		if dest != nil && resp.Result != nil {
			if err := json.Unmarshal(resp.Result, dest); err != nil {
				return fmt.Errorf("yanıt parse hatası: %w", err)
			}
		}
		return nil
	}

	return fmt.Errorf("telegram isteği başarısız (%d deneme): %w", maxRetries, lastErr)
}

func (c *Client) doRequest(ctx context.Context, method string, params map[string]any) (*apiResponse, error) {
	var bodyReader io.Reader
	var contentType string

	if params != nil {
		body, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("istek JSON hatası: %w", err)
		}
		bodyReader = bytes.NewReader(body)
		contentType = "application/json"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL(method), bodyReader)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP hatası: %w", err)
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(httpResp.Body, 10*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("yanıt okuma hatası: %w", err)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("API yanıt parse hatası: %w", err)
	}
	return &apiResp, nil
}

// sendMedia uploads binary data (photo/document) via multipart form.
func (c *Client) sendMedia(ctx context.Context, apiMethod, fieldName, filename string, chatID int64, data []byte, caption string) error {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	writer.WriteField("chat_id", strconv.FormatInt(chatID, 10))
	if caption != "" {
		writer.WriteField("caption", caption)
		writer.WriteField("parse_mode", "HTML")
	}

	part, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		return err
	}
	if _, err := part.Write(data); err != nil {
		return err
	}
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL(apiMethod), &body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("medya gönderme hatası: %w", err)
	}
	defer resp.Body.Close()

	var apiResp apiResponse
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	json.Unmarshal(respBody, &apiResp)

	if !apiResp.OK {
		return fmt.Errorf("medya gönderme API hatası: %s", apiResp.Description)
	}
	return nil
}

// ---- UTF-8 safe text helpers ----

// EscapeMarkdown escapes MarkdownV2 special characters in a string.
// Always use this before inserting user data into MarkdownV2 messages.
func EscapeMarkdown(s string) string {
	specialChars := `\_*[]()~` + "`" + `>#+-=|{}.!`
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s { // range iterates UTF-8 runes
		if strings.ContainsRune(specialChars, r) {
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}

// SafeText ensures the string is valid UTF-8 and safe for Telegram.
// Replaces any invalid UTF-8 sequences with the replacement character.
func SafeText(s string) string {
	if !isValidUTF8(s) {
		// Replace invalid sequences
		var b strings.Builder
		for _, r := range s {
			if r == '\uFFFD' {
				b.WriteRune('?')
			} else {
				b.WriteRune(r)
			}
		}
		return b.String()
	}
	return s
}

func isValidUTF8(s string) bool {
	for _, r := range s {
		if r == '\uFFFD' {
			return false
		}
	}
	return true
}
