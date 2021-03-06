# financoj

## Description
A commandline personal finance management program. It has capabilities of tracking expenses, creating budgets, reporting, among other features.

It uses reliable sqlite3 database to store your data, so it's safe and accessible using other tools as well.

## Installation

### Version 1.0
If you want to get the working version 1.0 then download the c sources from:
```
https://github.com/zbroju/financoj/archive/Release_1.0.zip
```
or
```
https://github.com/zbroju/financoj/archive/Release_1.0.tar.gz
```
and compile:
```
make
```
Version 1.0 requires the following c libraries to work:
```
libconfig
libsqlite3
```
### Current development
Current MASTER branch represents an early development stage of version 2 (where I move sources from c language to golang) and should be used *only* for development.

Download source code on your disk:
```
git clone https://github.com/zbroju/financoj
```
## Documentation:
Type:
```
fin --help
```
to get help and all available options
It is worth to copy the file example.financojrc to your $HOME/.financojrc and edit it by putting your own settings.
```
Usage:
	fin COMMAND [object | reports] [OPTIONS]
	
COMMANDS: 
        -I, --init	init a new file. Requires -f (--file) option.
        -A, --add	add new <object> to file.
        -E, --edit	edit existing <object>. Requires -i (--id) option to indicate the object.        
        -D, --delete	delete existing <object>. Requires -i (--id) option to indicate the object.        
        -L, --list	list <objects>. You can apply filters for the <objects>.        
        -R, --report	show <report>. You can apply filters for the <report>.        
        -h, --help	show this help information.
        
OBJECTS: 
        a, account	object to manipulate accounts.        
        t, transaction	object to manipulate transactions.        
        m, main-category	object to manipulate main categories.        
        j, currency	object to manipulate currencies.        
        c, category	object to manipulate categories.        
        b, budget	object to manipulate budgets.
        
REPORTS: 
        ab, accounts-balance	object to show report of accounts balances.        
        as, assets-summary	object to show report of assets summary.        
        tb, transactions-balance	object to show report of transactions balances.       
        cb, categories-balance	object to show report of categories balances.        
        mcb, main-categories-balance	object to show report of main categories balances.        
        bc, budget-categories	object to show report of budget for categories.        
        bmc, budget-main-categories	object to show report of budget for main categories.        
        nv, net-value	object to show report of net value.

OPTIONS: 
        -f, --file	full path to data file.        
        -i, --id	id for identifying particular object.        
        -n, --name	name of an object (account, main category & category).        
        -s, --description	description of a transaction.        
        -b, --bank	bank name where a given account is maintained.        
        -j, --currency	currency.        
        -k, --currency-to	currency against.        
        -a, --account	account name. It's enough to give part of the name as long as it allows to identify one account.        
        -c, --category	category name. It's enough to give part of the name as long as it allows to identify one category.        
        -m, --main-category	main category name. It's enough to give part of the name as long as it allows to identify one main category.
        -v, --value	value of a transaction, or exchange rate when working with currency.        
        -p, --account-type	account type. Allowed values are: t/transact (default), s/saving, p/property, i/investment, l/loan.        
        -o, --main-category-type	main category type. Allowed values are: c/cost, t/transfer, i/income.        
        -d, --date	date. Required format is YYYY (for year), YYYY-MM (for year-month) and YYYY-MM-DD (for full date). Today by default.        
        --verbose	make the program verbose.
```

## License
GNU General Public License

## Author
Marcin 'Zbroju' Zbroiński
