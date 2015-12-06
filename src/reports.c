/*
  Written 2015 by Marcin 'Zbroju' Zbroinski.
  Use of this source code is governed by a GNU General Public License
  that can be found in the LICENSE file.
*/

#include "reports.h"
#include "common.h"
#include <sqlite3.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>

#define REP_SQL_SIZE 10000
#define REP_BUF_SIZE 200
#define REP_TITLE_SIZE 100

#define EMP_ON "\x1b[1m"
#define EMP_OFF "\x1b[0m"

/**
 * The function returns true if currencies for all transactions have their
 * exchange rates towards reporting currency and false otherwise.
 * All the missing currencies will be placed as a string into parameters
 * missing_currencies.
 * If the function returns false, rememeber to free the memorey allocated for
 * missing_currencies (free(missing_currencies)), otherwise memory leak occurs.
 * @param db database pointer
 * @param reporting_currency string with reporting currency
 * @param missing_currencies string containing comma-separated list of missing currencies.
 * @return true if all currencies toward reporting one exist, or false otherwise.
 */
static bool all_currencies_available(sqlite3* db, char* reporting_currency, char** missing_currencies);

/* Function definitions */

int accounts_balance(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_accounts_balance[REP_SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    int result = 0;
    int year, month, day;
    ACCOUNT_TYPE current_type = ACC_TYPE_UNSET;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    }

    // Prepare date of balance
    if (parameters.date[0] != NULL_STRING) {
        if (date_from_string(parameters.date, &year, &month, &day) != DT_FULL_DATE) {
            fprintf(stderr, MSG_WRONG_DATE_FULL, parameters.prog_name);
            return 1;
        }
    } else {
        if (date_from_string(parameters.date_default, &year, &month, &day) != DT_FULL_DATE) {
            fprintf(stderr, MSG_WRONG_DATE_FULL, parameters.prog_name);
            return 1;
        }
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        return 1;
    }

    // Prepare SQL query and perform database actions
    sprintf(sql_accounts_balance,
            "SELECT"
            " A.TYPE" // 0
            ", A.NAME" // 1
            ", sum(T.VALUE)" // 2
            ", A.CURRENCY" // 3
            " FROM"
            " ("
            " SELECT ACCOUNT_ID, VALUE FROM TRANSACTIONS WHERE YEAR<%d"
            " UNION ALL"
            " SELECT ACCOUNT_ID, VALUE FROM TRANSACTIONS WHERE YEAR=%d AND MONTH<%d"
            " UNION ALL"
            " SELECT ACCOUNT_ID, VALUE FROM TRANSACTIONS WHERE YEAR=%d AND MONTH=%d AND DAY<=%d"
            ") T"
            " INNER JOIN ACCOUNTS A"
            " ON T.ACCOUNT_ID=A.ACCOUNT_ID"
            " WHERE A.STATUS=%d"
            " GROUP BY A.TYPE, A.NAME"
            " ORDER BY A.TYPE, A.NAME"
            ";"
            , year
            , year, month
            , year, month, day
            , ITEM_STAT_OPEN);
    if (sqlite3_prepare_v2(db, sql_accounts_balance, REP_SQL_SIZE, &sqlStmt, NULL) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        return 1;
    }

    printf("Accounts balance on " FS_DATE  "\n", year, month, day);

    // Print transactions on standard output
    while ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
        if (current_type != sqlite3_column_int(sqlStmt, 0)) {
            current_type = sqlite3_column_int(sqlStmt, 0);
            printf("\n" FS_ATYPE "\n", account_type_text(current_type));
        }
        printf(FS_GAP FS_NAME FS_GAP FS_VALUE FS_GAPS FS_CUR "\n"
                , sqlite3_column_text(sqlStmt, 1)
                , sqlite3_column_double(sqlStmt, 2)
                , sqlite3_column_text(sqlStmt, 3));
    }

    if (rc != SQLITE_DONE) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // Close database file
    rc = sqlite3_finalize(sqlStmt);

    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // if succeded and verbose then inform about performed actions
    // No need for special message if verbose.

    return result;
}

int assets_summary(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_accounts_balance[REP_SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    int result = 0;
    int year, month, day;
    ACCOUNT_TYPE current_type = ACC_TYPE_UNSET;
    char* list_of_missing_currencies = NULL;
    float subtotal =0 , total = 0;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.currency[0] == NULL_STRING) {
        if (parameters.default_currency[0] != NULL_STRING) {
            strncpy(parameters.currency
                    , parameters.default_currency, PAR_CURRENCY_LEN - 1);
        } else {
            fprintf(stderr, MSG_MISSING_PAR_CURRENCY
                    , parameters.prog_name
                    , OPTION_CURRENCY_SHORT
                    , OPTION_CURRENCY_LONG);
            return 1;
        }
    }

    // Prepare date of summary
    if (parameters.date[0] != NULL_STRING) {
        if (date_from_string(parameters.date, &year, &month, &day) != DT_FULL_DATE) {
            fprintf(stderr, MSG_WRONG_DATE_FULL, parameters.prog_name);
            return 1;
        }
    } else {
        if (date_from_string(parameters.date_default, &year, &month, &day) != DT_FULL_DATE) {
            fprintf(stderr, MSG_WRONG_DATE_FULL, parameters.prog_name);
            return 1;
        }
    }

    // Prepare SQL query and perform database actions
    sprintf(sql_accounts_balance,
            "SELECT"
            " A.TYPE" // 0
            ", A.NAME" // 1
            //", sum(T.VALUE)" // 2
            ", sum(CASE a.CURRENCY WHEN '%s' then 1 ELSE R.EXCHANGE_RATE END * T.VALUE) as CAT_VALUE" // 2
            " FROM"
            " ("
            " SELECT ACCOUNT_ID, VALUE FROM TRANSACTIONS WHERE YEAR<%d"
            " UNION ALL"
            " SELECT ACCOUNT_ID, VALUE FROM TRANSACTIONS WHERE YEAR=%d AND MONTH<%d"
            " UNION ALL"
            " SELECT ACCOUNT_ID, VALUE FROM TRANSACTIONS WHERE YEAR=%d AND MONTH=%d AND DAY<=%d"
            ") T"
            " INNER JOIN ACCOUNTS A"
            " ON T.ACCOUNT_ID=A.ACCOUNT_ID"
            " LEFT JOIN (SELECT * FROM CURRENCIES WHERE CURRENCY_TO='%s') R ON A.CURRENCY=R.CURRENCY_FROM"
            " WHERE A.STATUS=%d"
            " GROUP BY A.TYPE, A.NAME"
            " ORDER BY A.TYPE, A.NAME"
            ";"
            , parameters.currency
            , year
            , year, month
            , year, month, day
            , parameters.currency
            , ITEM_STAT_OPEN);

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        return 1;
    }

    if (all_currencies_available(db, parameters.currency, &list_of_missing_currencies) == true) {
        if (sqlite3_prepare_v2(db, sql_accounts_balance, REP_SQL_SIZE, &sqlStmt, NULL) != SQLITE_OK) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
            return 1;
        }

        printf("Assets summary on " FS_DATE  ":\n", year, month, day);

        // Print data on standard output
        while ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
            if (current_type != sqlite3_column_int(sqlStmt, 0)) {
                if (current_type != ACC_TYPE_UNSET) {
                    printf(EMP_ON FS_ATYPE "  " FS_GAP FS_VALUE FS_GAPS FS_CUR EMP_OFF"\n"
                            , account_type_text(current_type)
                            , subtotal
                            , parameters.currency);
                }
                current_type = sqlite3_column_int(sqlStmt, 0);
                printf("\n" FS_ATYPE  "\n", account_type_text(current_type));
                total += subtotal;
                subtotal = 0;
            }
            printf(FS_GAP FS_NAME FS_GAP FS_VALUE FS_GAPS FS_CUR "\n"
                    , sqlite3_column_text(sqlStmt, 1)
                    , sqlite3_column_double(sqlStmt, 2)
                    , parameters.currency);
            subtotal += sqlite3_column_double(sqlStmt, 2);
        }
        printf(EMP_ON FS_ATYPE "  " FS_GAP FS_VALUE FS_GAPS FS_CUR EMP_OFF "\n"
                , account_type_text(current_type)
                , subtotal
                , parameters.currency);
        total += subtotal;

        printf("\n" EMP_ON FS_NAME "  " FS_GAP FS_VALUE FS_GAPS FS_CUR EMP_OFF "\n"
               , "Total:"
                , total
                , parameters.currency);

        if (rc != SQLITE_DONE) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
            result = 1;
        }
    } else {
        fprintf(stderr, MSG_MISSING_EXCHANGE_RATES, parameters.prog_name, list_of_missing_currencies);
        free(list_of_missing_currencies);
    }

    // Close database file
    rc = sqlite3_finalize(sqlStmt);

    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // if succeded and verbose then inform about performed actions
    // No need for special message if verbose.

    return result;
}

int budget_report_categories(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_final[REP_SQL_SIZE] = {NULL_STRING};
    char sql_categories_list[REP_SQL_SIZE] = {NULL_STRING};
    char sql_budget[REP_SQL_SIZE] = {'0'};
    char sql_transactions[REP_SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    int result = 0;
    int year, month;
    char title[REP_TITLE_SIZE];
    float value_budget = 0.0,
            value_actual = 0.0,
            subtotal_budget = 0.0,
            subtotal_actual = 0.0,
            total_budget = 0.0,
            total_actual = 0.0;
    char* list_of_missing_currencies = NULL;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.currency[0] == NULL_STRING) {
        if (parameters.default_currency[0] != NULL_STRING) {
            strncpy(parameters.currency
                    , parameters.default_currency, PAR_CURRENCY_LEN - 1);
        } else {
            fprintf(stderr, MSG_MISSING_PAR_CURRENCY
                    , parameters.prog_name
                    , OPTION_CURRENCY_SHORT
                    , OPTION_CURRENCY_LONG);
            return 1;
        }
    }

    // Prepare date of budget
    DATE_TYPE date_type;
    if (parameters.date[0] != NULL_STRING) {
        date_type = date_from_string(parameters.date, &year, &month, NULL);
        if (date_type == DT_NO_DATE || date_type == DT_FULL_DATE) {
            fprintf(stderr, MSG_WRONG_DATE_MONTH_OR_YEAR, parameters.prog_name);
            return 1;
        }
    } else {
        date_type = date_from_string(parameters.date_default, &year, &month, NULL);
        if (date_type == DT_NO_DATE) {
            fprintf(stderr, MSG_WRONG_DATE_MONTH_OR_YEAR, parameters.prog_name);
            return 1;
        } else {
            date_type = DT_MONTH;
        }
    }

    // Prepare sql queries
    if (date_type == DT_MONTH) {
        sprintf(title, "Budget report for %d-%02d:", year, month);
        sprintf(sql_categories_list, "(SELECT DISTINCT CATEGORY_ID"
                " FROM BUDGETS"
                " WHERE YEAR=%d AND MONTH=%d"
                " UNION "
                " SELECT DISTINCT CATEGORY_ID"
                " FROM TRANSACTIONS"
                " WHERE YEAR=%d AND MONTH=%d)"
                , year, month
                , year, month);
        sprintf(sql_budget, "(SELECT BUDGETS.CATEGORY_ID,"
                " (BUDGETS.VALUE * (CASE BUDGETS.CURRENCY WHEN '%s' THEN 1.0 ELSE BUDGETCURRENCIES.EXCHANGE_RATE END)) AS BUDGET_LIMIT"
                " FROM BUDGETS"
                " LEFT JOIN (SELECT * FROM CURRENCIES WHERE CURRENCY_TO='%s') AS BUDGETCURRENCIES"
                " ON BUDGETS.CURRENCY=BUDGETCURRENCIES.CURRENCY_FROM"
                " WHERE BUDGETS.YEAR=%d AND BUDGETS.MONTH=%d)"
                , parameters.currency
                , parameters.currency
                , year, month);
        sprintf(sql_transactions, "(SELECT CATEGORY_ID"
                ", sum(VALUE * (CASE ACCOUNTS.CURRENCY WHEN '%s' THEN 1.0 ELSE TRANCURRENCIES.EXCHANGE_RATE END)) as SPENT_VALUE"
                " FROM TRANSACTIONS"
                " INNER JOIN ACCOUNTS ON TRANSACTIONS.ACCOUNT_ID=ACCOUNTS.ACCOUNT_ID"
                " LEFT JOIN (SELECT * FROM CURRENCIES WHERE CURRENCY_TO='%s') AS TRANCURRENCIES"
                " ON ACCOUNTS.CURRENCY=TRANCURRENCIES.CURRENCY_FROM"
                " WHERE YEAR=%d AND MONTH=%d"
                " GROUP BY CATEGORY_ID"
                ")"
                , parameters.currency
                , parameters.currency
                , year, month);
        sprintf(sql_final, "SELECT"
                " MAIN_CATEGORIES.TYPE"
                ", MAIN_CATEGORIES.NAME"
                ", CATEGORIES.NAME"
                ", CASE WHEN BUDGETLIMITS.BUDGET_LIMIT IS NULL THEN 0.0 ELSE BUDGETLIMITS.BUDGET_LIMIT END AS BUDGET"
                ", CASE WHEN TRANS.SPENT_VALUE IS NULL THEN 0.0 ELSE TRANS.SPENT_VALUE END AS ACTUAL"
                " FROM %s AS CATEGORIESLIST"
                " LEFT JOIN %s AS BUDGETLIMITS"
                " ON CATEGORIESLIST.CATEGORY_ID=BUDGETLIMITS.CATEGORY_ID"
                " LEFT JOIN %s AS TRANS"
                " ON CATEGORIESLIST.CATEGORY_ID=TRANS.CATEGORY_ID"
                " INNER JOIN CATEGORIES ON CATEGORIESLIST.CATEGORY_ID=CATEGORIES.CATEGORY_ID"
                " INNER JOIN MAIN_CATEGORIES ON CATEGORIES.MAIN_CATEGORY_ID=MAIN_CATEGORIES.MAIN_CATEGORY_ID"
                " ORDER BY MAIN_CATEGORIES.TYPE DESC, MAIN_CATEGORIES.NAME ASC"
                ";"
                , sql_categories_list
                , sql_budget
                , sql_transactions);

    } else {
        sprintf(title, "Budget report for %d:", year);
        sprintf(sql_categories_list, "(SELECT DISTINCT CATEGORY_ID"
                " FROM BUDGETS"
                " WHERE YEAR=%d"
                " UNION "
                " SELECT DISTINCT CATEGORY_ID"
                " FROM TRANSACTIONS"
                " WHERE YEAR=%d)"
                , year
                , year);
        sprintf(sql_budget, "(SELECT BUDGETS.CATEGORY_ID,"
                " sum(BUDGETS.VALUE * (CASE BUDGETS.CURRENCY WHEN '%s' THEN 1.0 ELSE BUDGETCURRENCIES.EXCHANGE_RATE END)) AS BUDGET_LIMIT"
                " FROM BUDGETS"
                " LEFT JOIN (SELECT * FROM CURRENCIES WHERE CURRENCY_TO='%s') AS BUDGETCURRENCIES"
                " ON BUDGETS.CURRENCY=BUDGETCURRENCIES.CURRENCY_FROM"
                " WHERE BUDGETS.YEAR=%d"
                " GROUP BY BUDGETS.CATEGORY_ID)"
                , parameters.currency
                , parameters.currency
                , year);
        sprintf(sql_transactions, "(SELECT CATEGORY_ID"
                ", sum(VALUE * (CASE ACCOUNTS.CURRENCY WHEN '%s' THEN 1.0 ELSE TRANCURRENCIES.EXCHANGE_RATE END)) as SPENT_VALUE"
                " FROM TRANSACTIONS"
                " INNER JOIN ACCOUNTS ON TRANSACTIONS.ACCOUNT_ID=ACCOUNTS.ACCOUNT_ID"
                " LEFT JOIN (SELECT * FROM CURRENCIES WHERE CURRENCY_TO='%s') AS TRANCURRENCIES"
                " ON ACCOUNTS.CURRENCY=TRANCURRENCIES.CURRENCY_FROM"
                " WHERE YEAR=%d"
                " GROUP BY CATEGORY_ID"
                ")"
                , parameters.currency
                , parameters.currency
                , year);
        sprintf(sql_final, "SELECT"
                " MAIN_CATEGORIES.TYPE"
                ", MAIN_CATEGORIES.NAME"
                ", CATEGORIES.NAME"
                ", CASE WHEN BUDGETLIMITS.BUDGET_LIMIT IS NULL THEN 0.0 ELSE BUDGETLIMITS.BUDGET_LIMIT END AS BUDGET"
                ", CASE WHEN TRANS.SPENT_VALUE IS NULL THEN 0.0 ELSE TRANS.SPENT_VALUE END AS ACTUAL"
                " FROM %s AS CATEGORIESLIST"
                " LEFT JOIN %s AS BUDGETLIMITS"
                " ON CATEGORIESLIST.CATEGORY_ID=BUDGETLIMITS.CATEGORY_ID"
                " LEFT JOIN %s AS TRANS"
                " ON CATEGORIESLIST.CATEGORY_ID=TRANS.CATEGORY_ID"
                " INNER JOIN CATEGORIES ON CATEGORIESLIST.CATEGORY_ID=CATEGORIES.CATEGORY_ID"
                " INNER JOIN MAIN_CATEGORIES ON CATEGORIES.MAIN_CATEGORY_ID=MAIN_CATEGORIES.MAIN_CATEGORY_ID"
                " ORDER BY MAIN_CATEGORIES.TYPE DESC, MAIN_CATEGORIES.NAME ASC"
                ";"
                , sql_categories_list
                , sql_budget
                , sql_transactions);
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        return 1;
    }

    // Are all the currencies exchange rates available?
    if (all_currencies_available(db, parameters.currency, &list_of_missing_currencies) == true) {

        if (sqlite3_prepare_v2(db, sql_final, REP_SQL_SIZE, &sqlStmt, NULL) != SQLITE_OK) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
            return 1;
        }

        // Print report on standard output
        MAINCATEGORY_TYPE prev_maincategory_type = CAT_TYPE_NOTSET, cur_maincategory_type;
        int subtotal_flag = 0;
        printf("%s\n", title);
        while ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
            cur_maincategory_type = sqlite3_column_int(sqlStmt, 0);
            if (cur_maincategory_type == CAT_INCOME || cur_maincategory_type == CAT_COST) {
                if (prev_maincategory_type != cur_maincategory_type) {
                    if (subtotal_flag == 0) {
                        subtotal_flag = 1;
                    } else {
                        printf(EMP_ON FS_MTYPE "    " FS_GAP FS_NAME FS_GAP FS_VALUE FS_GAPS FS_CUR FS_GAP FS_VALUE FS_GAPS FS_CUR FS_GAP FS_VALUE FS_GAPS FS_CUR EMP_OFF "\n"
                               , maincategory_type_text(prev_maincategory_type)
                               , ""
                                , subtotal_budget
                                , parameters.currency
                                , subtotal_actual
                                , parameters.currency
                                , subtotal_actual - subtotal_budget
                                , parameters.currency
                                );
                    }
                    total_budget += subtotal_budget;
                    total_actual += subtotal_actual;
                    subtotal_budget = 0.0;
                    subtotal_actual = 0.0;
                    prev_maincategory_type = cur_maincategory_type;

                    printf("\n" FS_MTYPE "\n", maincategory_type_text(cur_maincategory_type));
                    printf(FS_GAP FS_NAME_T FS_GAP FS_NAME_T FS_GAP FS_VALUE_T FS_GAPS FS_CUR_T FS_GAP FS_VALUE_T FS_GAPS FS_CUR_T FS_GAP FS_VALUE_T FS_GAPS FS_CUR_T"\n",
                           "MAIN CAT.", "CATEGORY", "LIMIT", "CUR", "ACTUAL", "CUR", "DIFFERENCE", "CUR");
                }
                printf(FS_GAP FS_NAME, sqlite3_column_text(sqlStmt, 1));
                printf(FS_GAP FS_NAME, sqlite3_column_text(sqlStmt, 2));
                value_budget = sqlite3_column_double(sqlStmt, 3);
                printf(FS_GAP FS_VALUE, value_budget);
                printf(FS_GAPS FS_CUR, parameters.currency);
                value_actual = sqlite3_column_double(sqlStmt, 4);
                printf(FS_GAP FS_VALUE, value_actual);
                printf(FS_GAPS FS_CUR, parameters.currency);
                printf(FS_GAP FS_VALUE, value_actual - value_budget);
                printf(FS_GAPS FS_CUR, parameters.currency);
                printf("\n");

                subtotal_budget += value_budget;
                subtotal_actual += value_actual;
            }
        }

        total_budget += subtotal_budget;
        total_actual += subtotal_actual;
        printf(EMP_ON FS_MTYPE "    " FS_GAP FS_NAME FS_GAP FS_VALUE FS_GAPS FS_CUR FS_GAP FS_VALUE FS_GAPS FS_CUR FS_GAP FS_VALUE FS_GAPS FS_CUR EMP_OFF"\n"
                , maincategory_type_text(prev_maincategory_type)
                , ""
                , subtotal_budget
                , parameters.currency
                , subtotal_actual
                , parameters.currency
                , subtotal_actual - subtotal_budget
                , parameters.currency
                );
        printf("\n" EMP_ON FS_MTYPE "    " FS_GAP FS_NAME FS_GAP FS_VALUE FS_GAPS FS_CUR FS_GAP FS_VALUE FS_GAPS FS_CUR FS_GAP FS_VALUE FS_GAPS FS_CUR EMP_OFF"\n"
                , "Total"
                , ""
                , total_budget
                , parameters.currency
                , total_actual
                , parameters.currency
                , total_actual - total_budget
                , parameters.currency
                );


        if (rc != SQLITE_DONE) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
            result = 1;
        }

        rc = sqlite3_finalize(sqlStmt);
    } else {
        fprintf(stderr, MSG_MISSING_EXCHANGE_RATES, parameters.prog_name, list_of_missing_currencies);
        free(list_of_missing_currencies);
    }
    // Close database file


    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // if succeded and verbose then inform about performed actions
    // No need for special message if verbose.

    return result;

}

int budget_report_maincategories(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_final[REP_SQL_SIZE] = {NULL_STRING};
    char sql_maincategories_list[REP_SQL_SIZE] = {NULL_STRING};
    char sql_budget[REP_SQL_SIZE] = {'0'};
    char sql_transactions[REP_SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    int result = 0;
    int year, month;
    char title[REP_TITLE_SIZE];
    float value_budget = 0.0,
            value_actual = 0.0,
            subtotal_budget = 0.0,
            subtotal_actual = 0.0,
            total_budget = 0.0,
            total_actual = 0.0;
    char* list_of_missing_currencies = NULL;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.currency[0] == NULL_STRING) {
        if (parameters.default_currency[0] != NULL_STRING) {
            strncpy(parameters.currency
                    , parameters.default_currency, PAR_CURRENCY_LEN - 1);
        } else {
            fprintf(stderr, MSG_MISSING_PAR_CURRENCY
                    , parameters.prog_name
                    , OPTION_CURRENCY_SHORT
                    , OPTION_CURRENCY_LONG);
            return 1;
        }
    }

    // Prepare date of budget
    DATE_TYPE date_type;
    if (parameters.date[0] != NULL_STRING) {
        date_type = date_from_string(parameters.date, &year, &month, NULL);
        if (date_type == DT_NO_DATE || date_type == DT_FULL_DATE) {
            fprintf(stderr, MSG_WRONG_DATE_MONTH_OR_YEAR, parameters.prog_name);
            return 1;
        }
    } else {
        date_type = date_from_string(parameters.date_default, &year, &month, NULL);
        if (date_type == DT_NO_DATE) {
            fprintf(stderr, MSG_WRONG_DATE_MONTH_OR_YEAR, parameters.prog_name);
            return 1;
        } else {
            date_type = DT_MONTH;
        }
    }

    // Prepare sql queries
    if (date_type == DT_MONTH) {
        sprintf(title, "Budget report for %d-%02d:", year, month);
        sprintf(sql_maincategories_list, "(SELECT DISTINCT MCL_C.MAIN_CATEGORY_ID AS MCL_MAIN_CATEGORY_ID FROM CATEGORIES AS MCL_C"
                " INNER JOIN (SELECT DISTINCT CATEGORY_ID"
                " FROM BUDGETS"
                " WHERE YEAR=%d AND MONTH=%d"
                " UNION "
                " SELECT DISTINCT CATEGORY_ID"
                " FROM TRANSACTIONS"
                " WHERE YEAR=%d AND MONTH=%d) AS MCL_UC ON MCL_UC.CATEGORY_ID=MCL_C.CATEGORY_ID)"
                , year, month
                , year, month);
        sprintf(sql_budget, "(SELECT B_C.MAIN_CATEGORY_ID, sum(B_BC.B_BUDGET_LIMIT)"
                " AS B_BUDGET_MAINCATS_LIMIT"
                " FROM (SELECT B_B.CATEGORY_ID,"
                " (B_B.VALUE * (CASE B_B.CURRENCY WHEN '%s' THEN 1.0 ELSE B_CURRENCIES.EXCHANGE_RATE END)) AS B_BUDGET_LIMIT"
                " FROM BUDGETS AS B_B"
                " LEFT JOIN (SELECT * FROM CURRENCIES WHERE CURRENCY_TO='%s') AS B_CURRENCIES"
                " ON B_B.CURRENCY=B_CURRENCIES.CURRENCY_FROM"
                " WHERE B_B.YEAR=%d AND B_B.MONTH=%d) AS B_BC"
                " INNER JOIN CATEGORIES AS B_C ON B_BC.CATEGORY_ID=B_C.CATEGORY_ID"
                " GROUP BY B_C.MAIN_CATEGORY_ID)"
                , parameters.currency
                , parameters.currency
                , year, month);
        sprintf(sql_transactions, "(SELECT T_C.MAIN_CATEGORY_ID"
                ", sum(VALUE * (CASE T_A.CURRENCY WHEN '%s' THEN 1.0 ELSE TRANCURRENCIES.EXCHANGE_RATE END)) as T_SPENT_VALUE"
                " FROM TRANSACTIONS AS T_T"
                " INNER JOIN ACCOUNTS AS T_A ON T_T.ACCOUNT_ID=T_A.ACCOUNT_ID"
                " LEFT JOIN (SELECT * FROM CURRENCIES WHERE CURRENCY_TO='%s') AS TRANCURRENCIES"
                " ON T_A.CURRENCY=TRANCURRENCIES.CURRENCY_FROM"
                " INNER JOIN CATEGORIES AS T_C ON T_C.CATEGORY_ID=T_T.CATEGORY_ID"
                " WHERE YEAR=%d AND MONTH=%d"
                " GROUP BY T_C.MAIN_CATEGORY_ID"
                ")"
                , parameters.currency
                , parameters.currency
                , year, month);
        sprintf(sql_final, "SELECT"
                " MAIN_CATEGORIES.TYPE"
                ", MAIN_CATEGORIES.NAME"
                ", CASE WHEN BUDGETLIMITS.B_BUDGET_MAINCATS_LIMIT IS NULL THEN 0.0 ELSE BUDGETLIMITS.B_BUDGET_MAINCATS_LIMIT END AS BUDGET"
                ", CASE WHEN TRANS.T_SPENT_VALUE IS NULL THEN 0.0 ELSE TRANS.T_SPENT_VALUE END AS ACTUAL"
                " FROM %s AS CATEGORIESLIST"
                " INNER JOIN MAIN_CATEGORIES ON CATEGORIESLIST.MCL_MAIN_CATEGORY_ID=MAIN_CATEGORIES.MAIN_CATEGORY_ID"
                " LEFT JOIN %s AS BUDGETLIMITS"
                " ON CATEGORIESLIST.MCL_MAIN_CATEGORY_ID=BUDGETLIMITS.MAIN_CATEGORY_ID"
                " LEFT JOIN %s AS TRANS"
                " ON CATEGORIESLIST.MCL_MAIN_CATEGORY_ID=TRANS.MAIN_CATEGORY_ID"
                " ORDER BY MAIN_CATEGORIES.TYPE DESC, MAIN_CATEGORIES.NAME ASC"
                ";"
                , sql_maincategories_list
                , sql_budget
                , sql_transactions);

    } else {
        sprintf(title, "Budget report for %d:", year);
        sprintf(sql_maincategories_list, "(SELECT DISTINCT MCL_C.MAIN_CATEGORY_ID AS MCL_MAIN_CATEGORY_ID FROM CATEGORIES AS MCL_C"
                " INNER JOIN (SELECT DISTINCT CATEGORY_ID"
                " FROM BUDGETS"
                " WHERE YEAR=%d"
                " UNION "
                " SELECT DISTINCT CATEGORY_ID"
                " FROM TRANSACTIONS"
                " WHERE YEAR=%d) AS MCL_UC ON MCL_UC.CATEGORY_ID=MCL_C.CATEGORY_ID)"
                , year
                , year);
        sprintf(sql_budget, "(SELECT B_C.MAIN_CATEGORY_ID, sum(B_BC.B_BUDGET_LIMIT)"
                " AS B_BUDGET_MAINCATS_LIMIT"
                " FROM (SELECT B_B.CATEGORY_ID,"
                " (B_B.VALUE * (CASE B_B.CURRENCY WHEN '%s' THEN 1.0 ELSE B_CURRENCIES.EXCHANGE_RATE END)) AS B_BUDGET_LIMIT"
                " FROM BUDGETS AS B_B"
                " LEFT JOIN (SELECT * FROM CURRENCIES WHERE CURRENCY_TO='%s') AS B_CURRENCIES"
                " ON B_B.CURRENCY=B_CURRENCIES.CURRENCY_FROM"
                " WHERE B_B.YEAR=%d) AS B_BC"
                " INNER JOIN CATEGORIES AS B_C ON B_BC.CATEGORY_ID=B_C.CATEGORY_ID"
                " GROUP BY B_C.MAIN_CATEGORY_ID)"
                , parameters.currency
                , parameters.currency
                , year);
        sprintf(sql_transactions, "(SELECT T_C.MAIN_CATEGORY_ID"
                ", sum(VALUE * (CASE T_A.CURRENCY WHEN '%s' THEN 1.0 ELSE TRANCURRENCIES.EXCHANGE_RATE END)) as T_SPENT_VALUE"
                " FROM TRANSACTIONS AS T_T"
                " INNER JOIN ACCOUNTS AS T_A ON T_T.ACCOUNT_ID=T_A.ACCOUNT_ID"
                " LEFT JOIN (SELECT * FROM CURRENCIES WHERE CURRENCY_TO='%s') AS TRANCURRENCIES"
                " ON T_A.CURRENCY=TRANCURRENCIES.CURRENCY_FROM"
                " INNER JOIN CATEGORIES AS T_C ON T_C.CATEGORY_ID=T_T.CATEGORY_ID"
                " WHERE YEAR=%d"
                " GROUP BY T_C.MAIN_CATEGORY_ID"
                ")"
                , parameters.currency
                , parameters.currency
                , year);
        sprintf(sql_final, "SELECT"
                " MAIN_CATEGORIES.TYPE"
                ", MAIN_CATEGORIES.NAME"
                ", CASE WHEN BUDGETLIMITS.B_BUDGET_MAINCATS_LIMIT IS NULL THEN 0.0 ELSE BUDGETLIMITS.B_BUDGET_MAINCATS_LIMIT END AS BUDGET"
                ", CASE WHEN TRANS.T_SPENT_VALUE IS NULL THEN 0.0 ELSE TRANS.T_SPENT_VALUE END AS ACTUAL"
                " FROM %s AS CATEGORIESLIST"
                " INNER JOIN MAIN_CATEGORIES ON CATEGORIESLIST.MCL_MAIN_CATEGORY_ID=MAIN_CATEGORIES.MAIN_CATEGORY_ID"
                " LEFT JOIN %s AS BUDGETLIMITS"
                " ON CATEGORIESLIST.MCL_MAIN_CATEGORY_ID=BUDGETLIMITS.MAIN_CATEGORY_ID"
                " LEFT JOIN %s AS TRANS"
                " ON CATEGORIESLIST.MCL_MAIN_CATEGORY_ID=TRANS.MAIN_CATEGORY_ID"
                " ORDER BY MAIN_CATEGORIES.TYPE DESC, MAIN_CATEGORIES.NAME ASC"
                ";"
                , sql_maincategories_list
                , sql_budget
                , sql_transactions);
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        return 1;
    }

    // Are all the currencies exchange rates available?
    if (all_currencies_available(db, parameters.currency, &list_of_missing_currencies) == true) {
        if (sqlite3_prepare_v2(db, sql_final, REP_SQL_SIZE, &sqlStmt, NULL) != SQLITE_OK) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
            return 1;
        }

        // Print report on standard output
        MAINCATEGORY_TYPE prev_maincategory_type = CAT_TYPE_NOTSET, cur_maincategory_type;
        int subtotal_flag = 0;
        printf("%s\n", title);
        while ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
            cur_maincategory_type = sqlite3_column_int(sqlStmt, 0);
            if (cur_maincategory_type == CAT_INCOME || cur_maincategory_type == CAT_COST) {
                if (prev_maincategory_type != cur_maincategory_type) {
                    if (subtotal_flag == 0) {
                        subtotal_flag = 1;
                    } else {
                        printf(EMP_ON FS_MTYPE "    " FS_GAP FS_VALUE FS_GAPS FS_CUR FS_GAP FS_VALUE FS_GAPS FS_CUR FS_GAP FS_VALUE FS_GAPS FS_CUR EMP_OFF "\n"
                                , maincategory_type_text(prev_maincategory_type)
                                , subtotal_budget
                                , parameters.currency
                                , subtotal_actual
                                , parameters.currency
                                , subtotal_actual - subtotal_budget
                                , parameters.currency
                                );
                    }
                    total_budget += subtotal_budget;
                    total_actual += subtotal_actual;
                    subtotal_budget = 0.0;
                    subtotal_actual = 0.0;
                    prev_maincategory_type = cur_maincategory_type;

                    printf("\n" FS_MTYPE "\n", maincategory_type_text(cur_maincategory_type));
                    printf(FS_GAP FS_NAME_T FS_GAP FS_VALUE_T FS_GAPS FS_CUR_T FS_GAP FS_VALUE_T FS_GAPS FS_CUR_T FS_GAP FS_VALUE_T FS_GAPS FS_CUR_T"\n",
                           "MAIN CAT.", "LIMIT", "CUR", "ACTUAL", "CUR", "DIFFERENCE", "CUR");
                }
                printf(FS_GAP FS_NAME, sqlite3_column_text(sqlStmt, 1));
                value_budget = sqlite3_column_double(sqlStmt, 2);
                printf(FS_GAP FS_VALUE, value_budget);
                printf(FS_GAPS FS_CUR, parameters.currency);
                value_actual = sqlite3_column_double(sqlStmt, 3);
                printf(FS_GAP FS_VALUE, value_actual);
                printf(FS_GAPS FS_CUR, parameters.currency);
                printf(FS_GAP FS_VALUE, value_actual - value_budget);
                printf(FS_GAPS FS_CUR, parameters.currency);
                printf("\n");

                subtotal_budget += value_budget;
                subtotal_actual += value_actual;
            }
        }

        total_budget += subtotal_budget;
        total_actual += subtotal_actual;
        printf(EMP_ON FS_MTYPE "    " FS_GAP FS_VALUE FS_GAPS FS_CUR FS_GAP FS_VALUE FS_GAPS FS_CUR FS_GAP FS_VALUE FS_GAPS FS_CUR EMP_OFF "\n"
                , maincategory_type_text(prev_maincategory_type)
                , subtotal_budget
                , parameters.currency
                , subtotal_actual
                , parameters.currency
                , subtotal_actual - subtotal_budget
                , parameters.currency
                );
        printf("\n" EMP_ON FS_MTYPE "    " FS_GAP FS_VALUE FS_GAPS FS_CUR FS_GAP FS_VALUE FS_GAPS FS_CUR FS_GAP FS_VALUE FS_GAPS FS_CUR EMP_OFF "\n"
                , "Total"
                , total_budget
                , parameters.currency
                , total_actual
                , parameters.currency
                , total_actual - total_budget
                , parameters.currency
                );


        if (rc != SQLITE_DONE) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
            result = 1;
        }

        rc = sqlite3_finalize(sqlStmt);
    } else {
        fprintf(stderr, MSG_MISSING_EXCHANGE_RATES, parameters.prog_name, list_of_missing_currencies);
        free(list_of_missing_currencies);
    }
    // Close database file


    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // if succeded and verbose then inform about performed actions
    // No need for special message if verbose.

    return result;

}

int categories_balance(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_categories_balance[REP_SQL_SIZE] = {NULL_STRING};
    char sql_transactions_subquery[REP_SQL_SIZE] = {NULL_STRING};
    char sql_currencies_subquery[REP_SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    int result = 0;
    int year, month, day;
    char balance_title[REP_TITLE_SIZE];
    float maincategory_balance = 0.0, total_balance = 0.0;
    MAINCATEGORY_TYPE maincategory_type = CAT_TYPE_NOTSET;
    char* list_of_missing_currencies = NULL;

    // Make sure necessary date have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.currency[0] == NULL_STRING) {
        if (parameters.default_currency[0] != NULL_STRING) {
            strncpy(parameters.currency
                    , parameters.default_currency
                    , PAR_CURRENCY_LEN - 1);
        } else {
            fprintf(stderr, MSG_MISSING_PAR_CURRENCY
                    , parameters.prog_name
                    , OPTION_CURRENCY_SHORT
                    , OPTION_CURRENCY_LONG);
            return 1;
        }
    }

    // Prepare (filter) date of balance
    DATE_TYPE date_type;
    if (parameters.date[0] != NULL_STRING) {
        if ((date_type = date_from_string(parameters.date, &year, &month, &day)) == DT_NO_DATE) {
            fprintf(stderr, MSG_WRONG_DATE_FULL, parameters.prog_name);
            return 1;
        } else {
            switch (date_type) {
            case DT_FULL_DATE:
                sprintf(sql_transactions_subquery,
                        "(SELECT * FROM TRANSACTIONS WHERE YEAR<%d"
                        " UNION ALL"
                        " SELECT * FROM TRANSACTIONS WHERE YEAR=%d AND MONTH<%d"
                        " UNION ALL"
                        " SELECT * FROM TRANSACTIONS WHERE YEAR=%d AND MONTH=%d AND DAY<=%d)"
                        , year
                        , year, month
                        , year, month, day
                        );
                sprintf(balance_title, "Category summary up to: %d-%02d-%02d", year, month, day);
                break;
            case DT_YEAR:
                sprintf(sql_transactions_subquery,
                        "(SELECT * FROM TRANSACTIONS WHERE YEAR=%d)"
                        , year);
                sprintf(balance_title, "Category summary during year %d", year);
                break;
            case DT_MONTH:
                sprintf(sql_transactions_subquery,
                        "(SELECT * FROM TRANSACTIONS WHERE YEAR=%d AND MONTH=%d)"
                        , year, month);
                sprintf(balance_title, "Category during month %d-%02d", year, month);
                break;
            case DT_NO_DATE:
                 break;
            }
        }
    } else {
        if ((date_type = date_from_string(parameters.date_default, &year, &month, &day)) == DT_NO_DATE) {
            fprintf(stderr, MSG_WRONG_DATE_FULL, parameters.prog_name);
            return 1;
        } else {
            sprintf(sql_transactions_subquery,
                    "(SELECT * FROM TRANSACTIONS WHERE YEAR<%d"
                    " UNION ALL"
                    " SELECT * FROM TRANSACTIONS WHERE YEAR=%d AND MONTH<%d"
                    " UNION ALL"
                    " SELECT * FROM TRANSACTIONS WHERE YEAR=%d AND MONTH=%d AND DAY<=%d)"
                    , year
                    , year, month
                    , year, month, day
                    );
            sprintf(balance_title, "Category summary up to: %d-%02d-%02d", year, month, day);
        }
    }

    // Prepare subquery for currencies
    sprintf(sql_currencies_subquery, "(SELECT * FROM CURRENCIES WHERE CURRENCY_TO='%s')",
            parameters.currency);

    // Prepare final query
    sprintf(sql_categories_balance,
            "SELECT"
            " mc.TYPE"
            ", mc.NAME"
            ", c.NAME"
            ", sum(CASE a.CURRENCY WHEN '%s' then 1 ELSE r.EXCHANGE_RATE END * t.VALUE) as CAT_VALUE"
            " FROM %s t"
            " INNER JOIN ACCOUNTS a ON t.ACCOUNT_ID=a.ACCOUNT_ID"
            " LEFT JOIN %s r ON a.CURRENCY=r.CURRENCY_FROM"
            " INNER JOIN CATEGORIES c ON t.CATEGORY_ID=c.CATEGORY_ID"
            " INNER JOIN MAIN_CATEGORIES mc ON c.MAIN_CATEGORY_ID=mc.MAIN_CATEGORY_ID"
            " GROUP BY mc.TYPE, mc.NAME, c.NAME"
            " ORDER BY mc.TYPE DESC, mc.NAME ASC, c.NAME ASC;"
            , parameters.currency
            , sql_transactions_subquery
            , sql_currencies_subquery);

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        return 1;
    }

    // Are all the currency exchange rates available?
    if (all_currencies_available(db, parameters.currency, &list_of_missing_currencies) == true) {
        if (sqlite3_prepare_v2(db, sql_categories_balance, REP_SQL_SIZE, &sqlStmt, NULL) != SQLITE_OK) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
            return 1;
        }

        // Print transactions on standard output
        printf("%s\n", balance_title);
        int subtotal_flag = 0;
        while ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
            if (maincategory_type != sqlite3_column_int(sqlStmt, 0)) {
                if (subtotal_flag == 0) {
                    subtotal_flag = 1;
                } else {
                    printf(EMP_ON FS_MTYPE "    " FS_GAP FS_NAME FS_GAP FS_VALUE FS_GAPS FS_CUR EMP_OFF "\n"
                            , maincategory_type_text(maincategory_type)
                           , ""
                            , maincategory_balance
                            , parameters.currency
                            );
                }
                maincategory_balance = 0.0;
                maincategory_type = sqlite3_column_int(sqlStmt, 0);
                printf("\n" FS_MTYPE "\n", maincategory_type_text(maincategory_type));
                printf(FS_GAP FS_NAME_T FS_GAP FS_NAME_T FS_GAP FS_VALUE_T FS_GAPS FS_CUR_T "\n"
                        , "MAIN CAT."
                        , "CATEGORY"
                        , "VALUE"
                        , "CUR");
            }
            printf(FS_GAP FS_NAME, sqlite3_column_text(sqlStmt, 1));
            printf(FS_GAP FS_NAME, sqlite3_column_text(sqlStmt, 2));
            printf(FS_GAP FS_VALUE, sqlite3_column_double(sqlStmt, 3));
            printf(FS_GAPS FS_CUR, parameters.currency);
            printf("\n");

            maincategory_balance += sqlite3_column_double(sqlStmt, 3);
            total_balance += sqlite3_column_double(sqlStmt, 3);
        }
        printf(EMP_ON FS_MTYPE "    " FS_GAP FS_NAME FS_GAP FS_VALUE FS_GAPS FS_CUR EMP_OFF "\n"
                , maincategory_type_text(maincategory_type)
               , ""
                , maincategory_balance
                , parameters.currency
                );
        printf("\n" EMP_ON FS_MTYPE "    " FS_GAP FS_NAME FS_GAP FS_VALUE FS_GAPS FS_CUR EMP_OFF "\n"
                , "Total"
               , ""
                , total_balance
                , parameters.currency
                );

        if (rc != SQLITE_DONE) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
            result = 1;
        }

        rc = sqlite3_finalize(sqlStmt);
    } else {
        fprintf(stderr, MSG_MISSING_EXCHANGE_RATES, parameters.prog_name, list_of_missing_currencies);
        free(list_of_missing_currencies);
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // if succeded and verbose then inform about performed actions
    // No need for special message if verbose.

    return result;
}

int maincategories_balance(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_categories_balance[REP_SQL_SIZE] = {NULL_STRING};
    char sql_transactions_subquery[REP_SQL_SIZE] = {NULL_STRING};
    char sql_currencies_subquery[REP_SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    int result = 0;
    int year, month, day;
    char balance_title[REP_TITLE_SIZE];
    float maincategory_balance = 0.0, total_balance = 0.0;
    MAINCATEGORY_TYPE maincategory_type = CAT_TYPE_NOTSET;
    char* list_of_missing_currencies = NULL;

    // Make sure necessary date have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.currency[0] == NULL_STRING) {
        if (parameters.default_currency[0] != NULL_STRING) {
            strncpy(parameters.currency
                    , parameters.default_currency
                    , PAR_CURRENCY_LEN - 1);
        } else {
            fprintf(stderr, MSG_MISSING_PAR_CURRENCY
                    , parameters.prog_name
                    , OPTION_CURRENCY_SHORT
                    , OPTION_CURRENCY_LONG);
            return 1;
        }
    }

    // Prepare (filter) date of balance
    DATE_TYPE date_type;
    if (parameters.date[0] != NULL_STRING) {
        if ((date_type = date_from_string(parameters.date, &year, &month, &day)) == DT_NO_DATE) {
            fprintf(stderr, MSG_WRONG_DATE_FULL, parameters.prog_name);
            return 1;
        } else {
            switch (date_type) {
            case DT_FULL_DATE:
                sprintf(sql_transactions_subquery,
                        "(SELECT * FROM TRANSACTIONS WHERE YEAR<%d"
                        " UNION ALL"
                        " SELECT * FROM TRANSACTIONS WHERE YEAR=%d AND MONTH<%d"
                        " UNION ALL"
                        " SELECT * FROM TRANSACTIONS WHERE YEAR=%d AND MONTH=%d AND DAY<=%d)"
                        , year
                        , year, month
                        , year, month, day
                        );
                sprintf(balance_title, "Category summary up to: %d-%02d-%02d", year, month, day);
                break;
            case DT_YEAR:
                sprintf(sql_transactions_subquery,
                        "(SELECT * FROM TRANSACTIONS WHERE YEAR=%d)"
                        , year);
                sprintf(balance_title, "Category summary during year %d", year);
                break;
            case DT_MONTH:
                sprintf(sql_transactions_subquery,
                        "(SELECT * FROM TRANSACTIONS WHERE YEAR=%d AND MONTH=%d)"
                        , year, month);
                sprintf(balance_title, "Category during month %d-%02d", year, month);
                break;
            case DT_NO_DATE:
                 break;
            }
        }
    } else {
        if ((date_type = date_from_string(parameters.date_default, &year, &month, &day)) == DT_NO_DATE) {
            fprintf(stderr, MSG_WRONG_DATE_FULL, parameters.prog_name);
            return 1;
        } else {
            sprintf(sql_transactions_subquery,
                    "(SELECT * FROM TRANSACTIONS WHERE YEAR<%d"
                    " UNION ALL"
                    " SELECT * FROM TRANSACTIONS WHERE YEAR=%d AND MONTH<%d"
                    " UNION ALL"
                    " SELECT * FROM TRANSACTIONS WHERE YEAR=%d AND MONTH=%d AND DAY<=%d)"
                    , year
                    , year, month
                    , year, month, day
                    );
            sprintf(balance_title, "Category summary up to: %d-%02d-%02d", year, month, day);
        }
    }

    // Prepare subquery for currencies
    sprintf(sql_currencies_subquery, "(SELECT * FROM CURRENCIES WHERE CURRENCY_TO='%s')",
            parameters.currency);

    // Prepare final query
    sprintf(sql_categories_balance,
            "SELECT"
            " mc.TYPE" // 0
            ", mc.NAME" // 1
            ", sum(CASE a.CURRENCY WHEN '%s' then 1 ELSE r.EXCHANGE_RATE END * t.VALUE) as CAT_VALUE" // 2
            " FROM %s t"
            " INNER JOIN ACCOUNTS a ON t.ACCOUNT_ID=a.ACCOUNT_ID"
            " LEFT JOIN %s r ON a.CURRENCY=r.CURRENCY_FROM"
            " INNER JOIN CATEGORIES c ON t.CATEGORY_ID=c.CATEGORY_ID"
            " INNER JOIN MAIN_CATEGORIES mc ON c.MAIN_CATEGORY_ID=mc.MAIN_CATEGORY_ID"
            " GROUP BY mc.TYPE, mc.NAME"
            " ORDER BY mc.TYPE DESC, mc.NAME ASC;"
            , parameters.currency
            , sql_transactions_subquery
            , sql_currencies_subquery);
    printf("\n%s\n\n", sql_categories_balance);
    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        return 1;
    }

    // Are all the currency exchange rates available?
    if (all_currencies_available(db, parameters.currency, &list_of_missing_currencies) == true) {

        if (sqlite3_prepare_v2(db, sql_categories_balance, REP_SQL_SIZE, &sqlStmt, NULL) != SQLITE_OK) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
            return 1;
        }

        // Print transactions on standard output
        printf("%s\n", balance_title);
        int subtotal_flag = 0;
        while ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
            if (maincategory_type != sqlite3_column_int(sqlStmt, 0)) {
                if (subtotal_flag == 0) {
                    subtotal_flag = 1;
                } else {
                    printf(EMP_ON FS_MTYPE "    " FS_GAP FS_VALUE FS_GAPS FS_CUR EMP_OFF "\n"
                            , maincategory_type_text(maincategory_type)
                            , maincategory_balance
                            , parameters.currency
                            );
                }
                maincategory_balance = 0.0;
                maincategory_type = sqlite3_column_int(sqlStmt, 0);
                printf("\n" FS_MTYPE  "\n", maincategory_type_text(maincategory_type));
                printf(FS_GAP FS_NAME_T FS_GAP FS_VALUE_T FS_GAPS FS_CUR_T "\n"
                        , "MAIN CAT."
                        , "VALUE"
                        , "CUR");
            }
            printf(FS_GAP FS_NAME, sqlite3_column_text(sqlStmt, 1));
            printf(FS_GAP FS_VALUE, sqlite3_column_double(sqlStmt, 2));
            printf(FS_GAPS FS_CUR, parameters.currency);
            printf("\n");

            maincategory_balance += sqlite3_column_double(sqlStmt, 2);
            total_balance += sqlite3_column_double(sqlStmt, 2);
        }
        printf(EMP_ON FS_MTYPE "    " FS_GAP FS_VALUE FS_GAPS FS_CUR EMP_OFF "\n"
                , maincategory_type_text(maincategory_type)
                , maincategory_balance
                , parameters.currency
                );
        printf("\n" EMP_ON FS_MTYPE "    " FS_GAP FS_VALUE FS_GAPS FS_CUR EMP_OFF "\n"
                , "Total"
                , total_balance
                , parameters.currency
                );

        if (rc != SQLITE_DONE) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
            result = 1;
        }

        rc = sqlite3_finalize(sqlStmt);
    } else {
        fprintf(stderr, MSG_MISSING_EXCHANGE_RATES, parameters.prog_name, list_of_missing_currencies);
        free(list_of_missing_currencies);
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // if succeded and verbose then inform about performed actions
    // No need for special message if verbose.

    return result;
}

int transactions_balance(PARAMETERS parameters)
{

    sqlite3 *db;
    char sql_transactions_balance[REP_SQL_SIZE] = {NULL_STRING};
    char sql_transactions_subquery[REP_SQL_SIZE] = {NULL_STRING};
    char sql_maincategory_subquery[REP_SQL_SIZE] = {NULL_STRING};
    char sql_category_subquery[REP_SQL_SIZE] = {NULL_STRING};
    char sql_account_subquery[REP_SQL_SIZE] = {NULL_STRING};
    char sql_currencies_subquery[REP_SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    int result = 0;
    int year, month, day;
    char balance_title[100];
    float balance_value = 0;
    char* list_of_missing_currencies = NULL;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.currency[0] == NULL_STRING) {
        if (parameters.default_currency[0] != NULL_STRING) {
            strncpy(parameters.currency
                    , parameters.default_currency, PAR_CURRENCY_LEN - 1);
        } else {
            fprintf(stderr, MSG_MISSING_PAR_CURRENCY
                    , parameters.prog_name
                    , OPTION_CURRENCY_SHORT
                    , OPTION_CURRENCY_LONG);
            return 1;
        }
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        return 1;
    }

    // Are all the currency exchange rates available?
    if (all_currencies_available(db, parameters.currency, &list_of_missing_currencies) == true) {

        // Prepare (filter) date of balance
        DATE_TYPE date_type;
        if (parameters.date[0] != NULL_STRING) {
            if ((date_type = date_from_string(parameters.date, &year, &month, &day)) == DT_NO_DATE) {
                fprintf(stderr, MSG_WRONG_DATE_FULL, parameters.prog_name);
                return 1;
            } else {
                switch (date_type) {
                case DT_FULL_DATE:
                    sprintf(sql_transactions_subquery,
                            "(SELECT * FROM TRANSACTIONS WHERE YEAR<%d"
                            " UNION ALL"
                            " SELECT * FROM TRANSACTIONS WHERE YEAR=%d AND MONTH<%d"
                            " UNION ALL"
                            " SELECT * FROM TRANSACTIONS WHERE YEAR=%d AND MONTH=%d AND DAY<=%d)"
                            , year
                            , year, month
                            , year, month, day
                            );
                    sprintf(balance_title, "Transactions up to: %d-%02d-%02d", year, month, day);
                    break;
                case DT_YEAR:
                    sprintf(sql_transactions_subquery,
                            "(SELECT * FROM TRANSACTIONS WHERE YEAR=%d)"
                            , year);
                    sprintf(balance_title, "Transactions during year %d", year);
                    break;
                case DT_MONTH:
                    sprintf(sql_transactions_subquery,
                            "(SELECT * FROM TRANSACTIONS WHERE YEAR=%d AND MONTH=%d)"
                            , year, month);
                    sprintf(balance_title, "Transactions during month %d-%02d", year, month);
                    break;
                case DT_NO_DATE:
                     break;
                }
            }
        } else {
            if ((date_type = date_from_string(parameters.date_default, &year, &month, &day)) == DT_NO_DATE) {
                fprintf(stderr, MSG_WRONG_DATE_FULL, parameters.prog_name);
                return 1;
            } else {
                sprintf(sql_transactions_subquery,
                        "(SELECT * FROM TRANSACTIONS WHERE YEAR<%d"
                        " UNION ALL"
                        " SELECT * FROM TRANSACTIONS WHERE YEAR=%d AND MONTH<%d"
                        " UNION ALL"
                        " SELECT * FROM TRANSACTIONS WHERE YEAR=%d AND MONTH=%d AND DAY<=%d)"
                        , year
                        , year, month
                        , year, month, day
                        );
                sprintf(balance_title, "Transactions up to: %d-%02d-%02d", year, month, day);
            }
        }

        // Prepare filter of main category
        long maincategory_id;
        if (parameters.maincat_name[0] == NULL_STRING) {
            sprintf(sql_maincategory_subquery, "(SELECT * FROM MAIN_CATEGORIES)");
        } else {
            maincategory_id = maincategory_id_for_name(db, parameters.maincat_name);
            if (maincategory_id == -1) {
                fprintf(stderr, MSG_MAINCATEGORY_NOT_FOUND
                        , parameters.prog_name, parameters.maincat_name);
                sqlite3_close(db);
                return 1;
            } else {
                sprintf(sql_maincategory_subquery,
                        "(SELECT * FROM MAIN_CATEGORIES WHERE MAIN_CATEGORY_ID=%li)",
                        maincategory_id);
            }
        }

        // Prepare filter of category
        long category_id;
        if (parameters.cat_name[0] == NULL_STRING) {
            sprintf(sql_category_subquery, "(SELECT * FROM CATEGORIES)");
        } else {
            category_id = category_id_for_name(db, parameters.cat_name);
            if (category_id == -1) {
                fprintf(stderr, MSG_CATEGORY_NOT_FOUND
                        , parameters.prog_name, parameters.cat_name);
                sqlite3_close(db);
                return 1;
            } else {
                sprintf(sql_category_subquery,
                        "(SELECT * FROM CATEGORIES WHERE CATEGORY_ID=%li)",
                        category_id);
            }
        }

        // Prepare filter of account
        long account_id;
        if (parameters.acc_name[0] == NULL_STRING) {
            sprintf(sql_account_subquery, "(SELECT * FROM ACCOUNTS)");
        } else {
            account_id = account_id_for_name(db, parameters.acc_name);
            if (account_id == -1) {
                fprintf(stderr, MSG_ACCOUNT_NOT_FOUND
                        , parameters.prog_name, parameters.acc_name);
                sqlite3_close(db);
                return 1;
            } else {
                sprintf(sql_account_subquery,
                        "(SELECT * FROM ACCOUNTS WHERE ACCOUNT_ID=%li)",
                        account_id);
            }
        }

        // Prepare subquery for currencies
        sprintf(sql_currencies_subquery, "(SELECT * FROM CURRENCIES WHERE CURRENCY_TO='%s')",
                parameters.currency);


        // Prepare final query
        sprintf(sql_transactions_balance,
                "SELECT"
                " u.NAME" // 0
                ", c.NAME" // 1
                ", a.NAME" // 2
                ", a.CURRENCY" // 3
                ", t.YEAR" // 4
                ", t.MONTH" // 5
                ", t.DAY" // 6
                ", t.DESCRIPTION" // 7
                ", CASE a.CURRENCY WHEN '%s' then 1 ELSE r.EXCHANGE_RATE END * t.VALUE as VALUE" // 8
                " FROM %s t"
                " INNER JOIN %s a ON t.ACCOUNT_ID=a.ACCOUNT_ID"
                " LEFT JOIN %s r ON a.CURRENCY=r.CURRENCY_FROM"
                " INNER JOIN %s c ON t.CATEGORY_ID=c.CATEGORY_ID"
                " INNER JOIN %s u ON c.MAIN_CATEGORY_ID=u.MAIN_CATEGORY_ID"
                " ORDER BY t.YEAR, t.MONTH, t.DAY"
                ";"
                , parameters.currency
                , sql_transactions_subquery
                , sql_account_subquery
                , sql_currencies_subquery
                , sql_category_subquery
                , sql_maincategory_subquery);



        if (sqlite3_prepare_v2(db, sql_transactions_balance, REP_SQL_SIZE, &sqlStmt, NULL) != SQLITE_OK) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
            return 1;
        }

        // Print transactions on standard output
        printf("%s\n", balance_title);
        printf("\n" FS_DATE_T FS_GAP FS_NAME_T FS_GAP FS_NAME_T FS_GAP FS_NAME_T FS_GAP FS_VALUE_T FS_GAPS FS_CUR_T FS_GAP FS_DESC_T "\n", "DATE", "MAIN CAT.", "CATEGORY", "ACCOUNT", "VALUE", "CUR", "DESCRIPTION");
        while ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
            printf(FS_DATE, sqlite3_column_int(sqlStmt, 4)
                    , sqlite3_column_int(sqlStmt, 5)
                    , sqlite3_column_int(sqlStmt, 6));
            printf(FS_GAP FS_NAME, sqlite3_column_text(sqlStmt, 0));
            printf(FS_GAP FS_NAME, sqlite3_column_text(sqlStmt, 1));
            printf(FS_GAP FS_NAME, sqlite3_column_text(sqlStmt, 2));
            printf(FS_GAP FS_VALUE, sqlite3_column_double(sqlStmt, 8));
            printf(FS_GAPS FS_CUR, parameters.currency);
            printf(FS_GAP FS_DESC, sqlite3_column_text(sqlStmt, 7));
            printf("\n");

            balance_value += sqlite3_column_double(sqlStmt, 8);
        }
        printf("\nTotal: " FS_VALUE FS_GAPS FS_CUR "\n", balance_value, parameters.currency);

        if (rc != SQLITE_DONE) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
            result = 1;
        }

        rc = sqlite3_finalize(sqlStmt);
    } else {
        fprintf(stderr, MSG_MISSING_EXCHANGE_RATES, parameters.prog_name, list_of_missing_currencies);
        free(list_of_missing_currencies);
    }

    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // if succeded and verbose then inform about performed actions
    // No need for special message if verbose.

    return result;
}

int net_value(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_net_value[REP_SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    int result = 0;
    char* list_of_missing_currencies = NULL;
    float net_value = 0;


    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.currency[0] == NULL_STRING) {
        if (parameters.default_currency[0] != NULL_STRING) {
            strncpy(parameters.currency
                    , parameters.default_currency, PAR_CURRENCY_LEN - 1);
        } else {
            fprintf(stderr, MSG_MISSING_PAR_CURRENCY
                    , parameters.prog_name
                    , OPTION_CURRENCY_SHORT
                    , OPTION_CURRENCY_LONG);
            return 1;
        }
    }

    // Prepare sql queries
    sprintf(sql_net_value, "SELECT"
            " YEAR"
            ",MONTH"
            ",sum(VALUE * (CASE A.CURRENCY WHEN '%s' THEN 1.0 ELSE TC.EXCHANGE_RATE END)) as NET_VALUE"
            " FROM TRANSACTIONS T"
            " INNER JOIN ACCOUNTS A ON T.ACCOUNT_ID=A.ACCOUNT_ID"
            " LEFT JOIN (SELECT * FROM CURRENCIES WHERE CURRENCY_TO='%s') AS TC"
            " ON A.CURRENCY=TC.CURRENCY_FROM"
            " GROUP BY YEAR, MONTH"
            " ORDER BY YEAR, MONTH;"
            , parameters.currency
            , parameters.currency
            );

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        return 1;
    }

    // Are all the currencies exchange rates available?
    if (all_currencies_available(db, parameters.currency, &list_of_missing_currencies) == true) {
        if (sqlite3_prepare_v2(db, sql_net_value, REP_SQL_SIZE, &sqlStmt, NULL) != SQLITE_OK) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
            return 1;
        }

        // Print report on standard output
        printf(FS_MONTH_T FS_GAP FS_VALUE_T FS_GAPS FS_CUR_T "\n", "PERIOD", "NET VALUE", "CUR");
        while ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
            net_value += sqlite3_column_double(sqlStmt, 2);
            printf(FS_MONTH FS_GAP FS_VALUE FS_GAPS FS_CUR "\n"
                    , sqlite3_column_int(sqlStmt, 0)
                    , sqlite3_column_int(sqlStmt, 1)
                    , net_value
                    , parameters.currency);
        }
        if (rc != SQLITE_DONE) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
            result = 1;
        }

        rc = sqlite3_finalize(sqlStmt);
    } else {
        fprintf(stderr, MSG_MISSING_EXCHANGE_RATES, parameters.prog_name, list_of_missing_currencies);
        free(list_of_missing_currencies);
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // if succeded and verbose then inform about performed actions
    // No need for special message if verbose.

    return result;
}


/* Local Functions */

#define CUR_ITEM_LEN 10

static bool all_currencies_available(sqlite3* db, char* reporting_currency, char** missing_currencies)
{
    char sql_query[REP_SQL_SIZE] = {NULL_STRING};
    char currencies_item[CUR_ITEM_LEN] = {NULL_STRING};
    char* tmp_ptr = NULL;
    sqlite3_stmt *sqlStmt_loc;
    int rc;
    bool result = true;
    int amount_of_items = 0;

    sprintf(sql_query, "SELECT CURRENCY, '%s' FROM"
            " (SELECT DISTINCT A.CURRENCY FROM TRANSACTIONS AS T"
            " INNER JOIN ACCOUNTS AS A ON T.ACCOUNT_ID=A.ACCOUNT_ID"
            " WHERE A.CURRENCY<>'%s') AS UC LEFT JOIN"
            " (SELECT CURRENCY_FROM FROM CURRENCIES WHERE CURRENCY_TO='%s') AS AC"
            " ON UC.CURRENCY=AC.CURRENCY_FROM WHERE AC.CURRENCY_FROM IS NULL;"
            , reporting_currency
            , reporting_currency
            , reporting_currency
            );

    if (sqlite3_prepare_v2(db, sql_query, REP_SQL_SIZE, &sqlStmt_loc, NULL) == SQLITE_OK) {
        while ((rc = sqlite3_step(sqlStmt_loc)) == SQLITE_ROW) {
            result = false;
            if ((tmp_ptr = realloc(*missing_currencies, (++amount_of_items) * CUR_ITEM_LEN * sizeof (char))) != NULL) {
                if (*missing_currencies == NULL) {
                    *missing_currencies = tmp_ptr;
                    sprintf(*missing_currencies, "%s-%s"
                            , sqlite3_column_text(sqlStmt_loc, 0)
                            , sqlite3_column_text(sqlStmt_loc, 1));
                } else {
                    *missing_currencies = tmp_ptr;
                    sprintf(currencies_item, ", %s-%s"
                            , sqlite3_column_text(sqlStmt_loc, 0)
                            , sqlite3_column_text(sqlStmt_loc, 1));
                    strncat(*missing_currencies, currencies_item, CUR_ITEM_LEN);
                }
            }
        }
    }

    // Clean
    rc = sqlite3_finalize(sqlStmt_loc);

    return result;
}
#undef CUR_ITEM_LEN

//TODO: net value historically
