// Written 2016 by Marcin 'Zbroju' Zbroinski.
// Use of this source code is governed by a GNU General Public License
// that can be found in the LICENSE file.

package engine

// SQL_REPORT_BUDGET_CATEGORIES_MONTHLY is SQL string to get budget values vs actual transactions value on category granularity
// for given month.
//
// Parameters
// 1 - reporting currency (string)
// 2 - year (string)
// 3 - month (string)
// 4 - year (int)
// 5 - month (int)
// 6 - Main Category Type Transfer (int)
// 7- reporting currency (string)
// 8 - reporting currency (string)
// 9 - year (int)
// 10 - month (int)
// 11 - reporting currency (string)
// 12 - reporting currency (string)
// 13 - year (string)
// 14 - month (string)
const SQL_REPORT_BUDGET_CATEGORIES_MONTHLY string = `
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
    , ? as currency
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
            ,currency
        from
            budgets
            inner join (select currency_from, exchange_rate from currencies where currency_to=? union select ?, 1) cur on budgets.currency=cur.currency_from
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
            inner join (select currency_from, exchange_rate from currencies where currency_to=? union select ?, 1) cur on a.currency=cur.currency_from
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
// 1 - reporting currency (string)
// 2 - year (string)
// 3 - year (int)
// 4 - Main Category Type Transfer (int)
// 5 - reporting currency (string)
// 6 - reporting currency (string)
// 7 - reporting currency (string)
// 8 - year (int)
// 9 - reporting currency (string)
// 10 - reporting currency (string)
// 11 - year (string)
const SQL_REPORT_BUDGET_CATEGORIES_YEARLY string = `
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
    , ? as currency
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
            ,?
        from
            budgets
            inner join (select currency_from, exchange_rate from currencies where currency_to=? union select ?, 1) cur on budgets.currency=cur.currency_from
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
            inner join (select currency_from, exchange_rate from currencies where currency_to=? union select ?, 1) cur on a.currency=cur.currency_from
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
