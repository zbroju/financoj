// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package engine

// SQL_REPORT_BUDGET_CATEGORIES_MONTHLY is SQL string to get budget values vs actual transactions value on category granularity
// for given month.
//
// Parameters
// 1 - year (string)
// 2 - month (string)
// 3 - year (int)
// 4 - month (int)
// 5 - Main Category Type Transfer (int)
// 6- reporting currency (string)
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

// SQL_REPORT_BUDGET_CATEGORIES_YEARLY is SQL string to get budget values vs actual transactions value on category granularity
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
