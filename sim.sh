#!/bin/bash

getifi() {
	local out=$(ip route show $1) # running the command with given subnet/cidr argument
	local ifi=$(cut -d " " -f3 <<< "$out")
	if [ "$ifi" = "" ]; then
		fail -1 "No such network"
	fi

	echo "$ifi"
}

fail() {
	echo "$2" >&2
	killall nc sim.sh > /dev/null 2>&1
	kill -9 $$ > /dev/null 2>&1
	exit $1
}

usage () {
	# tab=$(printf "\t")
	cat << E0F >&2 # cat into stderr till E0F
 Usage:
	./sim.sh alloc <ip subnet/cidr> [first digit] [number of addrs]
	Allocates IP addrs.
	Example: ./sim.sh alloc 192.168.10.0/24 1 10
	./sim.sh free
	Free the IP addrs.
	./sim.sh sim <file.conf>
	Start simulation, read config from file. File format:
	 portnumber header
	 portnumber header
	 ....
	Example ./sim.sh sim sims.conf
	(config file: 
	 22 SSH-2.0-OpenSSH_9.6p1 Ubuntu-3ubuntu13.5
	 21 220 Please use https://mirror.accum.se/ whenever possible. )
E0F

fail 1 "" # now this will exit with code 1 and return empty string into stderr
}

net=$1
ifi=$(getifi $net)
echo "Device name: $ifi"
