build:
	mkdir -p builds
	go build -o builds/vpngater

install:
	mv builds/vpngater /usr/local/sbin/vpngate
	chmod +x /usr/local/sbin/vpngate