// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package lib

// sqlReportTransactionsBalance is SQL string to get transactions values recalculated to one currency.
//
// Parameters
// 1 - currency_to (string)
// 2 - currency_to (string)
// 3 - date_from (string)
// 4 - date_from (string)
// 5 - noStringParamForSQL
// 6 - date_to (string)
// 7 - date_to (string)
// 8 - noStringParamForSQL
// 9 - account_id (int)
// 10 - account_id (int)
// 11 - noIntParamForSQL
// 12 - category_id (int)
// 13 - category_id (int)
// 14 - noIntParamForSQL
// 15 - main_category_id (int)
// 16 - main_category_id (int)
// 17 - noIntParamForSQL
// 18 - main_category_type (int)
// 19 - main_category_type (int)
// 20 - noIntParamForSQL
const sqlReportTransactionsBalance = `
select
    t.id
    ,t.date
    ,t.description
    ,t.value
    ,a.id
    ,a.name
    ,a.description
    ,a.institution
    ,a.currency
    ,a.type
    ,a.status
    ,c.id
    ,c.name
    ,c.status
    ,m.id
    ,m.name
    ,m.status
    ,mt.id
    ,mt.name
    ,mt.factor
    ,mt.factor * t.value * cur.exchange_rate as balance
from
    transactions t
    inner join accounts a on t.account_id=a.id
    inner join categories c on t.category_id=c.id
    inner join main_categories m on c.main_category_id=m.id
    inner join main_categories_types mt on m.type_id=mt.id
    inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on a.currency=cur.currency_from
where 1=1
    and (t.date>=? or ?=?)
    and (t.date<=? or ?=?)
    and (a.id=? or ?=?)
    and (c.id=? or ?=?)
    and (m.id=? or ?=?)
    and (t.description like ? or ?=?)
order by
    t.date
;
`

// sqlReportCategoriesBalance is SQL string to get categories values recalculated to one currency.
//
// Parameters
// 1 - MCTypeTransfer
// 2 - currency_to (string)
// 3 - currency_to (string)
// 4 - date_from (string)
// 5 - date_from (string)
// 6 - NoStringParamForSQL
// 7 - date_to (string)
// 8 - date_to (string)
// 9 - NoStringParamForSQL
// 10 - account_id (int)
// 11 - account_id (int)
// 12 - NoIntParamForSQL
// 13 - category_id (int)
// 14 - category_id (int)
// 15 - NoIntParamForSQL
// 16 - main_category_id (int)
// 17 - main_category_id (int)
// 18 - NoIntParamForSQL
const sqlReportCategoriesBalance string = `
select
    m.id
    ,m.name
    ,m.status
    ,mt.id
    ,mt.name
    ,mt.factor
    ,c.id
    ,c.name
    ,c.status
    ,sum(mt.factor * t.value * cur.exchange_rate) as balance
from
    transactions t
    inner join accounts a on t.account_id=a.id
    inner join categories c on t.category_id=c.id
    inner join main_categories m on c.main_category_id=m.id
    inner join (select * from main_categories_types where id<>?) mt on m.type_id=mt.id
    inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on a.currency=cur.currency_from
where 1=1
    and (t.date>=? or ?=?)
    and (t.date<=? or ?=?)
    and (a.id=? or ?=?)
    and (c.id=? or ?=?)
    and (m.id=? or ?=?)
group by
    m.id
    ,m.name
    ,m.status
    ,mt.id
    ,mt.name
    ,mt.factor
    ,c.id
    ,c.name
    ,c.status
order by
    mt.id desc
    ,m.name asc
    ,c.name asc
;
`

// sqlReportCategoriesBalance is SQL string to get categories values recalculated to one currency.
//
// Parameters
// 1 - MCTypeTransfer (int)
// 2 - Currency_to (string)
// 3 - Currency_to (string)
// 4 - date_from (string)
// 5 - date_from (string)
// 6 - NoStringParamForSQL
// 7 - date_to (string)
// 8 - date_to (string)
// 9 - NoStringParamForSQL
// 10 - account_id (int)
// 11 - account_id (int)
// 12 - NoIntParamForSQL
// 13 - main_category_id (int)
// 14 - main_category_id (int)
// 15 - NoIntParamForSQL
const sqlReportMainCategoriesBalance string = `
select
    m.id
    ,m.name
    ,m.status
    ,mt.id
    ,mt.name
    ,mt.factor
    ,sum(mt.factor * t.value * cur.exchange_rate) as balance
from
    transactions t
    inner join accounts a on t.account_id=a.id
    inner join categories c on t.category_id=c.id
    inner join main_categories m on c.main_category_id=m.id
    inner join (select * from main_categories_types where id<>?) mt on m.type_id=mt.id
    inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on a.currency=cur.currency_from
where 1=1
    and (t.date>=? or ?=?)
    and (t.date<=? or ?=?)
    and (a.id=? or ?=?)
    and (m.id=? or ?=?)
group by
    m.id
    ,m.name
    ,m.status
    ,mt.id
    ,mt.name
    ,mt.factor
order by
    mt.id desc
    ,m.name asc
;
`

// sqlReportAssetsSummary is SQL string to get assets values recalculated to one currency
//
// Parameters
// 1 - currency_to (string)
// 2 - currency_to (string)
// 3 - date_to (string)
// 4 - itemStatusClosed (int)
const sqlReportAssetsSummary string = `
select
    a.id
    ,a.name
    ,a.description
    ,a.institution
    ,a.currency
    ,a.type
    ,a.status
    ,sum(mt.factor * t.value * cur.exchange_rate) as balance
from
    transactions t
    inner join categories c on t.category_id=c.id
    inner join main_categories mc on c.main_category_id=mc.id
    inner join main_categories_types mt on mc.type_id=mt.id
    inner join accounts a on t.account_id=a.id
    inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on a.currency=cur.currency_from
where 1=1
    and t.date<=?
    and a.status<>?
group by
    a.id
    ,a.name
    ,a.description
    ,a.institution
    ,a.currency
    ,a.type
    ,a.status
order by
    a.type
    ,a.name
;
`

// sqlReportBudgetCategoriesMonthly is SQL string to get budget values vs actual transactions value on category granularity
// for given month.
//
// Parameters
// 1 - year (string)
// 2 - month (string)
// 3 - year (int)
// 4 - month (int)
// 5 - Main Category Type Transfer (int)
// 6 - reporting currency (string)
// 7 - reporting currency (string)
// 8 - year (int)
// 9 - month (int)
// 10 - reporting currency (string)
// 11 - reporting currency (string)
// 12 - year (string)
// 13 - month (string)
const sqlReportBudgetCategoriesMonthly string = `
select
    m.id
    , m.name
    , m.status
    , mct.id
    , mct.name
    , mct.factor
    , c.id
    , c.name
    , c.status
    , coalesce(b.budget,0.0) budgetLimit
    , coalesce(ta.actual,0.0) actualValue
    , coalesce(ta.actual,0.0) - coalesce(b.budget,0.0) as difference
from
    -- list of categories which have either budget or transactions
    (select
        category_id as id
    from
        transactions
    where
        strftime('%Y',date)=?
        and strftime('%m',date)=?
    union
    select
        category_id
    from
        budgets
    where
        year=?
        and month=?
) lc

    -- categories details
    inner join categories c on lc.id=c.id

    -- main categories details
    inner join main_categories m on c.main_category_id=m.id

    -- main categories types
    inner join (select * from main_categories_types where id<>?) mct on m.type_id=mct.id

    -- budget details
    left join (
        select
            category_id
            ,value * mt.factor * cur.exchange_rate as budget
        from
            budgets
            inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on budgets.currency=cur.currency_from
            inner join categories c on category_id=c.id
            inner join main_categories m on c.main_category_id=m.id
            inner join main_categories_types mt on m.type_id=mt.id
        where
            year=?
            and month=?
    ) b on lc.id=b.category_id

    -- actual transactions
    left join (
        select
            t.category_id
            ,sum(t.value * mt.factor * cur.exchange_rate) as actual
        from
           transactions t
            inner join accounts a on t.account_id=a.id
            inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on a.currency=cur.currency_from
            inner join categories c on t.category_id=c.id
            inner join main_categories m on c.main_category_id=m.id
            inner join main_categories_types mt on m.type_id=mt.id
        where
            strftime('%Y',t.date)=?
            and strftime('%m',t.date)=?
        group by
            category_id
    ) ta on lc.id=ta.category_id
order by
	mct.factor DESC
	, m.name
	, c.name
;
`

// sqlReportBudgetCategoriesYearly is SQL string to get budget values vs actual transactions value on category granularity
// for given year.
//
// Parameters
// 1 - year (string)
// 2 - year (int)
// 3 - Main Category Type Transfer (int)
// 4 - reporting currency (string)
// 5 - reporting currency (string)
// 6 - year (int)
// 7 - reporting currency (string)
// 8 - reporting currency (string)
// 9 - year (string)
const sqlReportBudgetCategoriesYearly string = `
select
    m.id
    , m.name
    , m.status
    , mct.id
    , mct.name
    , mct.factor
    , c.id
    , c.name
    , c.status
    , coalesce(b.budget,0.0) budgetLimit
    , coalesce(ta.actual,0.0) actualValue
    , coalesce(ta.actual,0.0) - coalesce(b.budget,0.0) as difference
from
    -- list of categories which have either budget or transactions
    (select
        category_id as id
    from
        transactions
    where
        strftime('%Y',date)=?
    union
    select
        category_id
    from
        budgets
    where
        year=?
) lc

    -- categories details
    inner join categories c on lc.id=c.id

    -- main categories details
    inner join main_categories m on c.main_category_id=m.id

    -- main categories types
    inner join (select * from main_categories_types where id<>?) mct on m.type_id=mct.id

    -- budget details
    left join (
        select
            category_id
            ,sum(value * mt.factor * cur.exchange_rate) as budget
        from
            budgets
            inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on budgets.currency=cur.currency_from
            inner join categories c on category_id=c.id
            inner join main_categories m on c.main_category_id=m.id
            inner join main_categories_types mt on m.type_id=mt.id
        where
            year=?
        group by
            category_id
    ) b on lc.id=b.category_id

    -- actual transactions
    left join (
        select
            t.category_id
            ,sum(t.value * mt.factor * cur.exchange_rate) as actual
        from
           transactions t
            inner join accounts a on t.account_id=a.id
            inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on a.currency=cur.currency_from
            inner join categories c on t.category_id=c.id
            inner join main_categories m on c.main_category_id=m.id
            inner join main_categories_types mt on m.type_id=mt.id
        where
            strftime('%Y',t.date)=?
        group by
            category_id
    ) ta on lc.id=ta.category_id
order by
	mct.factor DESC
	, m.name
	, c.name
;
`

// sqlReportBudgetMainCategoriesMonthly is SQL string to get budget values vs actual transactions value on main category granularity
// for given month.
//
// Parameters
// 1 - year (string)
// 2 - month (string)
// 3 - year (int)
// 4 - month (int)
// 5 - main category transfer id (int)
// 6 - currency_to (string)
// 7 - currency_to (string)
// 8 - year (int)
// 9 - month (int)
// 10 - year (string)
// 11 - month (string)
// 12 - currency_to (string)
// 13 - currency_to (string)
const sqlReportBudgetMainCategoriesMonthly string = `
select
    lmc.id
    , lmc.name
    , lmc.status
    , mct.id
    , mct.name
    , mct.factor
    , coalesce(b.budget,0.0) budgetLimit
    , coalesce(ta.actual,0.0) actualValue
    , coalesce(ta.actual,0.0) - coalesce(b.budget,0.0) as difference
from
    -- list of main categories which have either budget or transactions
    (select
    	mc.*
    from
        main_categories mc inner join categories c on mc.id=c.main_category_id inner join (select * from transactions where strftime('%Y',date)=? and strftime('%m',date)=?) t on c.id=t.category_id
    union
    select
        mb.*
    from
        main_categories mb inner join categories c on mb.id=c.main_category_id inner join (select * from budgets where year=? and month=?) b on c.id=b.category_id
) lmc

    -- main categories types
    inner join (select * from main_categories_types where id<>?) mct on lmc.type_id=mct.id

    -- budget details
    left join (
        select
            m.id
            ,sum(value * mt.factor * cur.exchange_rate) as budget
        from
            budgets
            inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on budgets.currency=cur.currency_from
            inner join categories c on category_id=c.id
            inner join main_categories m on c.main_category_id=m.id
            inner join main_categories_types mt on m.type_id=mt.id
        where
            year=?
            and month=?
        group by
            m.id
    ) b on lmc.id=b.id

    -- actual transactions
    left join (
        select
            m.id
            ,sum(t.value * mt.factor * cur.exchange_rate) as actual
        from
           (select * from transactions where strftime('%Y',date)=? and strftime('%m',date)=?) t
            inner join accounts a on t.account_id=a.id
            inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on a.currency=cur.currency_from
            inner join categories c on t.category_id=c.id
            inner join main_categories m on c.main_category_id=m.id
            inner join main_categories_types mt on m.type_id=mt.id
        group by
            m.id
    ) ta on lmc.id=ta.id
order by
	mct.factor DESC
	,lmc.name
;
`

// sqlReportBudgetMainCategoriesYearly is SQL string to get budget values vs actual transactions value on main category granularity
// for given year.
//
// Parameters
// 1 - year (string)
// 2 - year (int)
// 3 - main category type transfer (int)
// 4 - currency_to (string)
// 5 - currency_to (string)
// 6 - year (int)
// 7 - year (string)
// 8 - currency_to (string)
// 9 - currency_to (string)
const sqlReportBudgetMainCategoriesYearly string = `
select
    lmc.id
    , lmc.name
    , lmc.status
    , mct.id
    , mct.name
    , mct.factor
    , coalesce(b.budget,0.0) budgetLimit
    , coalesce(ta.actual,0.0) actualValue
    , coalesce(ta.actual,0.0) - coalesce(b.budget,0.0) as difference
from
    -- list of main categories which have either budget or transactions
    (select
    	mc.*
    from
        main_categories mc inner join categories c on mc.id=c.main_category_id inner join (select * from transactions where strftime('%Y',date)=?) t on c.id=t.category_id
    union
    select
        mb.*
    from
        main_categories mb inner join categories c on mb.id=c.main_category_id inner join (select * from budgets where year=?) b on c.id=b.category_id
) lmc

    -- main categories types
    inner join (select * from main_categories_types where id<>?) mct on lmc.type_id=mct.id

    -- budget details
    left join (
        select
            m.id
            ,sum(value * mt.factor * cur.exchange_rate) as budget
        from
            budgets
            inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on budgets.currency=cur.currency_from
            inner join categories c on category_id=c.id
            inner join main_categories m on c.main_category_id=m.id
            inner join main_categories_types mt on m.type_id=mt.id
        where
            year=?
        group by
            m.id
    ) b on lmc.id=b.id

    -- actual transactions
    left join (
        select
            m.id
            ,sum(t.value * mt.factor * cur.exchange_rate) as actual
        from
           (select * from transactions where strftime('%Y',date)=?) t
            inner join accounts a on t.account_id=a.id
            inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on a.currency=cur.currency_from
            inner join categories c on t.category_id=c.id
            inner join main_categories m on c.main_category_id=m.id
            inner join main_categories_types mt on m.type_id=mt.id
        group by
            m.id
    ) ta on lmc.id=ta.id
order by
	mct.factor DESC
	,lmc.name
;
`

// sqlReportMissingCurrenciesForTransactions is SQL string to get all currencies used in transactions where there is no exchange rate.
//
// Parameters:
// 1 - currency (string)
// 2 - currency (string)
// 3 - currency (string)
const sqlReportMissingCurrenciesForTransactions string = `
select currency || '-' || upper(?) as cur_to
from
(select distinct a.currency from transactions t inner join accounts a on t.account_id=a.id where a.currency<>upper(?)) uc
left join (select currency_from from currencies where currency_to=upper(?)) ac on uc.currency=ac.currency_from
where
ac.currency_from is null
;
`

// sqlReportNetValueMonthly is SQL string to get monthly balance of all transactions in order to build net value.
//
// Paramters:
// 1 - currency (string)
// 2 - currency (string)
// 3 - date_from (string)
// 4 - date_from (string)
// 5 - NoStringParamForSQL
// 6 - date_to (string)
// 7 - date_to (string)
// 8 - NoStringParamForSQL
const sqlReportNetValueMonthly string = `
select
    strftime('%Y',date) as y
    ,strftime('%m',date) as m
    ,sum(mt.factor * t.value * cur.exchange_rate) as balance
from
    transactions t
    inner join categories c on t.category_id=c.id
    inner join main_categories mc on c.main_category_id=mc.id
    inner join main_categories_types mt on mc.type_id=mt.id
    inner join accounts a on t.account_id=a.id
    inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on a.currency=cur.currency_from
where 1=1
    and (t.date>=? or ?=?)
    and (t.date<=? or ?=?)
group by
    y
    ,m
order by
    y
    ,m
;
`

// sqlReportMissingCurrenciesForBudgets is SQL string to get all currencies used in transactions where there is no exchange rate.
//
// Parameters:
// 1 - currency (string)
// 2 - currency (string)
// 3 - currency (string)
const sqlReportMissingCurrenciesForBudgets string = `
select currency || '-' || upper(?) as cur_to
from
	(select distinct currency from budgets where currency<>upper(?)) uc
	left join (select currency_from from currencies where currency_to=upper(?)) ac on uc.currency=ac.currency_from
where
	ac.currency_from is null
;
`

// sqlReportCategoriesBalanceMonthly is SQL string to get selected category balance over time (monthly).
//
// Parameters
// 1 - currency (string)
// 2 - currency (string)
// 3 - category_id (int)
// 4 - date_from (string)
// 5 - date_from (string)
// 6 - NoStringParamForSQL
// 7 - date_to (string)
// 8 - date_to (string)
// 9 - NoStringParamForSQL
const sqlReportCategoriesBalanceMonthly string = `
select
	strftime('%Y',date) as y
    	,strftime('%m',date) as m
	,sum(t.value * mt.factor * cur.exchange_rate) as balance
from
	transactions t
	inner join accounts a on t.account_id=a.id
	inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on a.currency=cur.currency_from
	inner join categories c on t.category_id=c.id
	inner join main_categories m on c.main_category_id=m.id
	inner join main_categories_types mt on m.type_id=mt.id
where 1=1
	and c.id=?
	and (t.date>=? or ?=?)
	and (t.date<=? or ?=?)
group by
	1
	,2
order by
	1
	,2
;
`

// sqlReportCategoriesBalanceYearly is SQL string to get selected category balance over time (yearly).
//
// Parameters
// 1 - currency (string)
// 2 - currency (string)
// 3 - category_id (int)
// 4 - date_from (string)
// 5 - date_from (string)
// 6 - NoStringParamForSQL
// 7 - date_to (string)
// 8 - date_to (string)
// 9 - NoStringParamForSQL
const sqlReportCategoriesBalanceYearly string = `
select
	strftime('%Y',date) as y
	,sum(t.value * mt.factor * cur.exchange_rate) as balance
from
	transactions t
	inner join accounts a on t.account_id=a.id
	inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on a.currency=cur.currency_from
	inner join categories c on t.category_id=c.id
	inner join main_categories m on c.main_category_id=m.id
	inner join main_categories_types mt on m.type_id=mt.id
where 1=1
	and c.id=?
	and (t.date>=? or ?=?)
	and (t.date<=? or ?=?)
group by
	1
order by
	1
;
`

// sqlReportMainCategoriesBalanceMonthly is SQL string to get selected main category balance over time (monthly).
//
// Parameters
// 1 - currency (string)
// 2 - currency (string)
// 3 - main_category_id (int)
// 4 - date_from (string)
// 5 - date_from (string)
// 6 - NoStringParamForSQL
// 7 - date_to (string)
// 8 - date_to (string)
// 9 - NoStringParamForSQL
const sqlReportMainCategoriesBalanceMonthly string = `
select
	strftime('%Y',date) as y
    	,strftime('%m',date) as m
	,sum(t.value * mt.factor * cur.exchange_rate) as balance
from
	transactions t
	inner join accounts a on t.account_id=a.id
	inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on a.currency=cur.currency_from
	inner join categories c on t.category_id=c.id
	inner join main_categories m on c.main_category_id=m.id
	inner join main_categories_types mt on m.type_id=mt.id
where 1=1
	and m.id=?
	and (t.date>=? or ?=?)
	and (t.date<=? or ?=?)
group by
	1
	,2
order by
	1
	,2
;
`

// sqlReportMainCategoriesBalanceYearly is SQL string to get selected main category balance over time (yearly).
//
// Parameters
// 1 - currency (string)
// 2 - currency (string)
// 3 - main_category_id (int)
// 4 - date_from (string)
// 5 - date_from (string)
// 6 - NoStringParamForSQL
// 7 - date_to (string)
// 8 - date_to (string)
// 9 - NoStringParamForSQL
const sqlReportMainCategoriesBalanceYearly string = `
select
	strftime('%Y',date) as y
	,sum(t.value * mt.factor * cur.exchange_rate) as balance
from
	transactions t
	inner join accounts a on t.account_id=a.id
	inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on a.currency=cur.currency_from
	inner join categories c on t.category_id=c.id
	inner join main_categories m on c.main_category_id=m.id
	inner join main_categories_types mt on m.type_id=mt.id
where 1=1
	and m.id=?
	and (t.date>=? or ?=?)
	and (t.date<=? or ?=?)
group by
	1
order by
	1
;
`

// sqlReportIncomeAndCostMonthly is SQL string to get sum of income and cost per month
//
// Parameters
// 1 - date_from (string)
// 2 - date_from (string)
// 3 - NoStringParamForSQL
// 4 - date_to (string)
// 5 - date_to (string
// 6 - NoStringParamForSQL
// 7 - currency (string)
// 8 - currency (string)
// 9 - MainCategoryTypeIncome (int)
// 10 - currency (string)
// 11 - currency (string)
// 12 - MainCategoryTypeCost (int)
const sqlReportIncomeAndCostMonthly string = `
-- PERIODS
select
	periods.year
	,periods.month
	,coalesce(income.balance, 0.0) as income
	,coalesce(cost.balance, 0.0) as cost
from
	(select
		strftime('%Y', t.date) as year
		,strftime('%m', t.date) as month
	from
		transactions t
	where 1=1
		and (t.date>=? or ?=?)
		and (t.date<=? or ?=?)
	group by 1,2
	order by 1,2) periods
	left join
	-- INCOME
	(select
		strftime('%Y',t.date) as year
		,strftime('%m',t.date) as month
		,sum(mt.factor * t.value * cur.exchange_rate) as balance
	from
		transactions t
		inner join accounts a on t.account_id=a.id
		inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on a.currency=cur.currency_from
		inner join categories c on t.category_id=c.id
		inner join main_categories m on c.main_category_id=m.id
		inner join main_categories_types mt on m.type_id=mt.id
	where
		mt.id=?
	group by
		1
		,2) income on periods.year=income.year and periods.month=income.month
	left join
	-- COST
	(select
		strftime('%Y',t.date) as year
		,strftime('%m',t.date) as month
		,sum(mt.factor * t.value * cur.exchange_rate) as balance
	from
		transactions t
		inner join accounts a on t.account_id=a.id
		inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on a.currency=cur.currency_from
		inner join categories c on t.category_id=c.id
		inner join main_categories m on c.main_category_id=m.id
		inner join main_categories_types mt on m.type_id=mt.id
	where
		mt.id=?
	group by
		1
		,2) cost on periods.year=cost.year and periods.month=cost.month
;
`

// sqlReportIncomeAndCostYearly is SQL string to get sum of income and cost per year
//
// Parameters
// 1 - date_from (string)
// 2 - date_from (string)
// 3 - NoStringParamForSQL
// 4 - date_to (string)
// 5 - date_to (string
// 6 - NoStringParamForSQL
// 7 - currency (string)
// 8 - currency (string)
// 9 - MainCategoryTypeIncome (int)
// 10 - currency (string)
// 11 - currency (string)
// 12 - MainCategoryTypeCost (int)
const sqlReportIncomeAndCostYearly string = `
-- PERIODS
select
	periods.year
	,coalesce(income.balance, 0.0) as income
	,coalesce(cost.balance, 0.0) as cost
from
	(select
		strftime('%Y', t.date) as year
	from
		transactions t
	where 1=1
		and (t.date>=? or ?=?)
		and (t.date<=? or ?=?)
	group by 1
	order by 1) periods
	left join
	-- INCOME
	(select
		strftime('%Y',t.date) as year
		,sum(mt.factor * t.value * cur.exchange_rate) as balance
	from
		transactions t
		inner join accounts a on t.account_id=a.id
		inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on a.currency=cur.currency_from
		inner join categories c on t.category_id=c.id
		inner join main_categories m on c.main_category_id=m.id
		inner join main_categories_types mt on m.type_id=mt.id
	where
		mt.id=?
	group by
		1) income on periods.year=income.year
	left join
	-- COST
	(select
		strftime('%Y',t.date) as year
		,sum(mt.factor * t.value * cur.exchange_rate) as balance
	from
		transactions t
		inner join accounts a on t.account_id=a.id
		inner join (select currency_from, exchange_rate from currencies where currency_to=upper(?) union select upper(?), 1) cur on a.currency=cur.currency_from
		inner join categories c on t.category_id=c.id
		inner join main_categories m on c.main_category_id=m.id
		inner join main_categories_types mt on m.type_id=mt.id
	where
		mt.id=?
	group by
		1) cost on periods.year=cost.year
;
`
