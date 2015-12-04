# Description and ideas
- two level categories
- one budget per year split by months
- based on sqlite
- only basic operations: income, expenses, budgets, credits, properties, net value
- assume working with one file only one year: archiving and smooth switch (initial balances) to new year

# Main functions

## Operational accounts
- list
- add
- edit
- close
- has name, description, initial balance, actvie/inactive, institution, currency
- different type: operational (related to budgets), credit, properties - related to net value.

## Categories
- list
- add
- edit
- remove, but keep the history (in DB make it inactive)
- has two levels: category:subcategory;
- has type: income, expense, transfer

## Transaction
- list
- add
- remove
- edit
- has account, date, type:income/expense/transfer, subcategory, value, description
- all transactions are budgetable!

## Currencies
- list
- add
- edit
- remove
- has: from date, to date, currency from, currency to, exchange rate

## Archiving
- archive transactions for date range in separate file
- for current file copy only end date balance from old one to opening balance in current

## Budget
- One budget per year split in months
- all categories shown: only values per month to be entered

## Reporting
- home page: accounts summary, budgets summary, yearly expenses/income per category in pie charts, 
- income/expenses per type & category (subtotals) for a date range, plus pie charts
- income/expenses per category per months in a year
- budget vs actual comparison: yearly, per month 
