# TO BE DONE
- implement formatting strings so that the reports are compact
- implement full strings for user input
- add signature to the database – the same as for gsqlitehandler and check it every time when opening the file
- write tutorial
- prepare README with installation instruction and dependencies
- prepare debian package

# COMPLETED
- create application
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
- add comments to example of configuration file
- make comparison for currencies insensitive to case letters
- add condition that an account cannot be closed if its balance is <> 0
- change name of the project to financoj

# DISCARDED
- add function to check if the given file is a mymoney database (RATIONALE: not necessary: application doesn't overwrite or change existing file)
- overview report (balances of all operational accounts, summaries of other accounts, current net value, summary of budget); (RATIONALE: not necessary, existing reports are enough)
- add condition to check if object exists in listings and reports (e.g. account); (RATIONALE: not necessary, if there is no objects to show, only header is printed, which is OK.)
- package it as deb & tar; (RATIONALE: not necessary until someone needs that.)
- deleting budget for the whole month/year; (RATIONALE: there is no issue with data size, it's better to keep control over it but forcing removing the budgets one by one.)
- copy budget to new month; (RATIONALE: not neccessary because it can be easily done with bash scripts.)
- add report pre-formatted for gnuplot (net value, categories balance); (RATIONALE: it doesn't make sense, because we need mainly bar and pie charts and pie charts are not very easily done in gnuplot.)
