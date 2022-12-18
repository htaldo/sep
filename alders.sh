#!/bin/bash
#https://stackoverflow.com/questions/339483/how-can-i-remove-the-first-line-of-a-text-file-using-bash-sed-script#339941

#CAMPANA
awk 'BEGIN { OFS = "\t" } { print $1, $2, $3 }' reparto >extracto
awk 'BEGIN { OFS = "\t" } { print $4, $5, $6 }' reparto | tac >refinado
cat extracto refinado | \
awk 'BEGIN { OFS = "\t"
             prev0 = "" }
     { print $0, prev0
       prev0=$0 }' | tail -n +2 >campana
#ALDERS
awk 'BEGIN { OFS = "\t" } { print $1, $2, $3, $1, $5, 1 - $1 - $5 }' reparto >interp
awk 'BEGIN { OFS = "\t" } { print $4, $5, $6, $1, $5, 1 - $1 - $5 }' reparto >>interp
awk 'BEGIN { OFS = "\t" } 
     { schnittp = $1 "\t" $5 "\t" 1 - $1 - $5 
       print schnittp, prevschnittp
       prevschnittp = schnittp }' reparto | tail -n +2 >alders
