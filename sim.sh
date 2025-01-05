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

fail 1 "$1" # now this will exit with code 1 and return empty string into stderr
}

alloc () {
	echo "Running alloc ..."
	net=$1
	first=$2
	n=$3

	if let "!n"; then
		let "n=10"
	fi
	if let "!first"; then
		let "first=1"
	fi
	if [ "$net" = "" ]; then
		usage "No network given"
	fi

	ifi=$(getifi $net)
	# getting the first 3 digits of an ip
	st=$(cut -d "." -f1-3 <<< "$net") # as in start
	own_ip=$(ip addr show "$ifi" |egrep 'inet '| awk '{print $2}'| cut -d '/' -f1)
	echo "Own IP: $own_ip"
	echo "Ifi: $ifi"
	if [ "$own_ip" = "" ]; then
		own_ip="300.300.300.300"
	fi

}

free() {
	echo "Running free ..."
}

sim() {
	echo "Running sim ..."
}

cmd=$1 # the first would be a command like alloc, free, sim
uid=$(id -u)
if let "$uid"; then # so if it is not 0 -> not root
	fail -1 "You must provide root access"
fi

case $cmd in 
	"alloc")
	alloc $2 $3 $4 # 2 = ip subnet; 3 = first dig; 4 = num of addrs
	;;
	"free")
	free
	;;
	"sim")
	sim $2 # $2 = confing file
	;;
esac
