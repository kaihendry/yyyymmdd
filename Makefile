APP = yyyymmdd
FN =  ~/Library/LaunchAgents/managedl.plist


install:
	cp managedl.plist $(FN)

list:
	launchctl list | grep $(APP)

load:
	launchctl load -w $(FN)

logs:
	tail -F /var/log/system.log | grep --line-buffered $(APP)
