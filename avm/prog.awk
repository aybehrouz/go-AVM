BEGIN {
	FS = "[.:, \t]+"
}

/0x/ {
	print $2 "\t" $5
}