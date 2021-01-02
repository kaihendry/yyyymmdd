Trying to organise my downloads into YYYY-MM-DD folders since my browsers don't

	~/.config/systemd/user/managedl.service

# Issues

I should have kept the files in ~/Downloads and not moved them to ~/dl because
files might be closed though they are still being operated upon, like .part
files. Then you get odd situations like:

	[hendry@t480s ~]$ ls -l ~/dl/2020-10-04
	total 27712
	-rw-r--r-- 1 hendry users        0 Oct  4 12:17 'openkiosk-68.3.0-2020-01-30-x86_64(1).tar.bz2'
	-rw------- 1 hendry users 27213824 Oct  4 12:20 'openkiosk-68.3.0-2020-01-30-x86_64(1).tar.bz2.part'
	-rw-r--r-- 1 hendry users        0 Oct  4 12:17 'openkiosk-68.3.0-2020-01-30-x86_64(2).tar.bz2'
	-rw------- 1 hendry users  1163264 Oct  4 12:17 'openkiosk-68.3.0-2020-01-30-x86_64(2).tar.bz2.part'
