#!/bin/sh -eu

GOT=$(cat)
WANT="$@"

if test "$GOT" = "$WANT"; then
  echo "Output matched assertion"
  exit 0
fi

echo "Output did not match assertion"

echo ">>>"
echo "$GOT"
echo ">>>"

echo "<<<"
echo "$WANT"
echo "<<<"

exit 1
