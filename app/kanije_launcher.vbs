Set WshShell = CreateObject("WScript.Shell")
WshShell.CurrentDirectory = "C:\\Users\\Yasin\\Downloads\\Kanije-Kalesi\\app"
WshShell.Run "pythonw ""C:\\Users\\Yasin\\Downloads\\Kanije-Kalesi\\app\\kanije.py"" start", 0, False
