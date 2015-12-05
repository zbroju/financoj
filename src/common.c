/*
  Written 2015 by Marcin 'Zbroju' Zbroinski.
  Use of this source code is governed by a GNU General Public License
  that can be found in the LICENSE file.
*/

#include "common.h"
#include <stdio.h>
#include <string.h>
#include <time.h>

#define COMMON_SQL_SIZE 1000

long account_id_for_name(sqlite3* db, char* acc_name)
{
    char sql_getId[COMMON_SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    long result = 0;
    int number_of_accounts = 0;

    // Prepare SQL query
    sprintf(sql_getId, "SELECT ACCOUNT_ID FROM ACCOUNTS WHERE NAME LIKE '%%%s%%' AND STATUS=%d;"
            , acc_name
            , ITEM_STAT_OPEN);

    if (sqlite3_prepare_v2(db, sql_getId, COMMON_SQL_SIZE, &sqlStmt, NULL) != SQLITE_OK) {
        result = -1;
    }

    // Count accounts
    while ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
        number_of_accounts++;
        result = sqlite3_column_int(sqlStmt, 0);
    }

    // Return id only if no error and there is only one account id for given name
    if ((rc != SQLITE_DONE) || (number_of_accounts != 1)) {
        result = -1;
    }

    // Clean
    rc = sqlite3_finalize(sqlStmt);

    return result;
}

ACCOUNT_TYPE account_type_id(char* account_type)
{
    ACCOUNT_TYPE result = ACC_TYPE_UNKNOWN;

    if (strncmp(account_type, "transact", OBJ_OR_TYPE_LEN) == 0
            || strncmp(account_type, "t", OBJ_OR_TYPE_LEN) == 0) {
        result = ACC_TRANSACTIONAL;
    } else if (strncmp(account_type, "saving", OBJ_OR_TYPE_LEN) == 0
            || strncmp(account_type, "s", OBJ_OR_TYPE_LEN) == 0) {
        result = ACC_SAVING;
    } else if (strncmp(account_type, "property", OBJ_OR_TYPE_LEN) == 0
            || strncmp(account_type, "p", OBJ_OR_TYPE_LEN) == 0) {
        result = ACC_PROPERTY;
    } else if (strncmp(account_type, "investment", OBJ_OR_TYPE_LEN) == 0
            || strncmp(account_type, "i", OBJ_OR_TYPE_LEN) == 0) {
        result == ACC_INVESTMENT;
    } else if (strncmp(account_type, "loan", OBJ_OR_TYPE_LEN) == 0
            || strncmp(account_type, "l", OBJ_OR_TYPE_LEN) == 0) {
        result = ACC_LOAN;
    }

    return result;
}

char* account_type_text(ACCOUNT_TYPE acc_type)
{
    switch (acc_type) {
    case ACC_TYPE_UNKNOWN:
        return "unknown";
    case ACC_TYPE_UNSET:
        return "unset";
    case ACC_TRANSACTIONAL:
        return "Operations";
    case ACC_SAVING:
        return "Savings";
    case ACC_PROPERTY:
        return "Investment";
    case ACC_LOAN:
        return "Loan";
    default:
        return "unknown";
    }
}

int category_factor_for_id(sqlite3* db, int cat_id)
{
    char sql_getType[COMMON_SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    int result = 0;
    MAINCATEGORY_TYPE main_category_type;

    // Prepare SQL query
    sprintf(sql_getType, "SELECT M.TYPE"
            " FROM CATEGORIES C INNER JOIN MAIN_CATEGORIES M"
            " ON C.MAIN_CATEGORY_ID=M.MAIN_CATEGORY_ID"
            " WHERE C.CATEGORY_ID=%d;"
            , cat_id);
    if (sqlite3_prepare_v2(db, sql_getType, COMMON_SQL_SIZE, &sqlStmt, NULL) != SQLITE_OK) {
        main_category_type = CAT_TYPE_UNKNOWN;
    }

    // Take type
    while ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
        main_category_type = sqlite3_column_int(sqlStmt, 0);
    }

    // Find the factor
    result = maincategory_type_factor(main_category_type);

    // Clean
    rc = sqlite3_finalize(sqlStmt);

    return result;
}

long category_id_for_name(sqlite3* db, char* cat_name)
{
    char sql_getId[COMMON_SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    long result = 0;
    int number_of_categories = 0;

    // Prepare SQL query
    sprintf(sql_getId, "SELECT CATEGORY_ID FROM CATEGORIES WHERE NAME LIKE '%%%s%%' AND STATUS=%d;", cat_name, ITEM_STAT_OPEN);

    if (sqlite3_prepare_v2(db, sql_getId, COMMON_SQL_SIZE, &sqlStmt, NULL) != SQLITE_OK) {
        result = -1;
    }

    // Count categories
    while ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
        number_of_categories++;
        result = sqlite3_column_int(sqlStmt, 0);
    }

    // Return id only if no error and there is only one category id for given name
    if ((rc != SQLITE_DONE) || (number_of_categories != 1)) {
        result = -1;
    }

    // Clean
    rc = sqlite3_finalize(sqlStmt);

    return result;
}

DATE_TYPE date_from_string(char* date_string, int *year_holder, int *month_holder, int *day_holder)
{
    DATE_TYPE result = DT_NO_DATE;
    int year = 0, month = 0, day = 0;

    switch (strlen(date_string)) {
    case 10:
        if (sscanf(date_string, "%d-%d-%d", &year, &month, &day) != EOF) {
            result = DT_FULL_DATE;
        }
        break;
    case 7:
        if (sscanf(date_string, "%d-%d", &year, &month) != EOF) {
            result = DT_MONTH;
        }
        break;
    case 4:
        if (sscanf(date_string, "%d", &year) != EOF) {
            result = DT_YEAR;
        }
        break;
    }

    if (result != DT_NO_DATE) {
        if (year_holder != NULL) {
            *year_holder = year;
        }
        if (month_holder != NULL) {
            *month_holder = month;
        }
        if (day_holder != NULL) {
            *day_holder = day;
        }
    }

    return result;
}

void get_today(char* date_holder)
{
    time_t time_now;
    struct tm *today;

    time_now = time(NULL);
    if (time_now != -1) {
        today = localtime(&time_now);
        strftime(date_holder, DATE_FULL_LEN, "%Y-%m-%d", today);
    } else {
        date_holder[0] = NULL_STRING;
    }
}

char* item_status_text(ITEM_STATUS item_status)
{
    switch (item_status) {
    case ITEM_STAT_CLOSED:
        return "Closed";
    case ITEM_STAT_OPEN:
        return "Open";
    default:
        return "Unknown";
    }
}

long maincategory_id_for_name(sqlite3* db, char* maincategory_name)
{
    char sql_getId[COMMON_SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    long result = 0;
    int number_of_maincategories = 0;

    // Prepare SQL query
    sprintf(sql_getId, "SELECT MAIN_CATEGORY_ID FROM MAIN_CATEGORIES WHERE NAME LIKE '%%%s%%' AND STATUS=%d;"
            , maincategory_name
            , ITEM_STAT_OPEN);
    if (sqlite3_prepare_v2(db, sql_getId, COMMON_SQL_SIZE, &sqlStmt, NULL) != SQLITE_OK) {
        result = -1;
    }

    // Count categories
    while ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
        number_of_maincategories++;
        result = sqlite3_column_int(sqlStmt, 0);
    }

    // Return id only if no error and there is only one category id for given name
    if ((rc != SQLITE_DONE) || (number_of_maincategories != 1)) {
        result = -1;
    }

    // Clean
    rc = sqlite3_finalize(sqlStmt);

    return result;
}

int maincategory_type_factor(MAINCATEGORY_TYPE maincategory_type)
{
    switch (maincategory_type) {
    case CAT_COST:
        return -1;
    case CAT_TRANSFER:
        return 1;
    case CAT_INCOME:
        return 1;
    default:
        return 0;
    }
}

MAINCATEGORY_TYPE maincategory_type_id(char* maincategory_type)
{
    MAINCATEGORY_TYPE result = CAT_TYPE_UNKNOWN;

    if (strncmp(maincategory_type, "cost", OBJ_OR_TYPE_LEN) == 0
            || strncmp(maincategory_type, "c", OBJ_OR_TYPE_LEN) == 0) {
        result = CAT_COST;
    } else if (strncmp(maincategory_type, "transfer", OBJ_OR_TYPE_LEN) == 0
            || strncmp(maincategory_type, "t", OBJ_OR_TYPE_LEN) == 0) {
        result = CAT_TRANSFER;
    } else if (strncmp(maincategory_type, "income", OBJ_OR_TYPE_LEN) == 0
            || strncmp(maincategory_type, "i", OBJ_OR_TYPE_LEN) == 0) {
        result = CAT_INCOME;
    }

    return result;
}

char* maincategory_type_text(MAINCATEGORY_TYPE maincategory_type)
{
    switch (maincategory_type) {
    case CAT_COST:
        return "Cost";
    case CAT_TRANSFER:
        return "Transfer";
    case CAT_INCOME:
        return "Income";
    default:
        return "unset";
    }
}

unsigned long transaction_category_id(sqlite3* db, unsigned long transaction_id)
{
    char sql_getCategoryID[COMMON_SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    int result = 0;


    // Prepare SQL query
    sprintf(sql_getCategoryID, "SELECT CATEGORY_ID FROM TRANSACTIONS WHERE TRANSACTION_ID=%d;",
            transaction_id);
    if (sqlite3_prepare_v2(db, sql_getCategoryID, COMMON_SQL_SIZE, &sqlStmt, NULL) != SQLITE_OK) {
        result = -1;
    }

    // Take type
    while ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
        result = sqlite3_column_int(sqlStmt, 0);
    }

    // Clean
    rc = sqlite3_finalize(sqlStmt);

    return result;
}
