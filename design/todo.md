# TO BE DONE
- add literals to definition of short options for getopt
- add help for each command separately
- add function to check if the given file is a mymoney database
- add condition to check if object exists in listings and reports (e.g. account)
- add report pre-formatted for gnuplot (net value, categories balance)
- add condition that an account cannot be closed if its balance is <> 0
- make comparison for currencies insensitive to case letters
- overview report (balances of all operational accounts, summaries of other accounts, current net value, summary of budget)
- tutorial step by step how to work
- package it as deb & tar
- deleting budget for the whole month/year
- copy budget to new month


# COMPLETED
- create applicatoin
- report with account balance in one currency (assets)
- secure the license is attached
- report: net value
- check if showing account balance excludes closed accounts
- check if showing assets balance excludes closed accounts
- check if getID account excludes closed items
- check if getID category excludes closed items
- check if getID main category excludes closed items
- change comments in all files to be consistent and aligned with license
- round values inserted/edited to db up to 2 digits after comma
- round exchange rate iserted/edited to db up to 4 digits after comma
- set string length for all object (account name, main category, category, description etc.),
- truncate all the strings before inserting to DB according to the set string lengths
- adjust reports so that respect string lengths and be more compact
