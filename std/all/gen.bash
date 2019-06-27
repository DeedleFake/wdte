#!/bin/bash

out=/dev/stdout
pkgname='all'

while [ -n "$1" ]; do
	case "$1" in
		-o)
			shift
			out="$1"
			;;

		-p)
			shift
			pkgname="$1"
			;;

		*)
			echo "Unrecognized argument: $1"
			exit 2
			;;
	esac

	shift
done

pkgs=$(cd .. && find . -type d -mindepth 1 | grep -v '\./all')

echo "package $pkgname" > $out
echo >> $out
echo '// Code generated automatically. DO NOT EDIT.' >> $out
echo >> $out
echo "import (" >> $out
for pkg in $pkgs; do
	echo "	_ \"github.com/DeedleFake/wdte/std/$(echo "$pkg" | cut -c3-)\""
done | sort >> $out
echo ")" >> $out
