/*
  Written 2015 by Marcin 'Zbroju' Zbroinski.
  Use of this source code is governed by a GNU General Public License
  that can be found in the LICENSE file.
*/

/* INCLUDES */
#include "operations.h"
#include "common.h"
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <sqlite3.h>
#include <time.h>
#include <stdbool.h>

#define SQL_SIZE 1000
#define BUF_SIZE 200

/**
 * Returns true if an account for given id exists, or false otherwise.
 * @param db sqlite3* database pointer
 * @param account_id account_id unsigned long with id of the account.
 * @return true or false
 */
static bool account_exists(sqlite3* db, const unsigned long account_id);

/**
 * Returns true if a budget for given date and category exists, or false otherwise.
 * @param db sqlite3* database pointer
 * @param category_id an unsigned long with category id of the budget
 * @param year an int with a year number
 * @param month an int with a month number
 * @return true or false
 */
static bool budget_exists(sqlite3* db, const unsigned long category_id, const int year, const int month);

/**
 * Returns true if a category for given id exists, or false otherwise.
 * @param db sqlite3* database pointer
 * @param category_id category_id unsigned long with id of the account.
 * @return true or false
 */
static bool category_exists(sqlite3* db, const unsigned long category_id);

/**
 * Returns true if an exchange rate for given pair of currencies exists, or false otherwise.
 * @param db sqlite3* database pointer
 * @param currency_from char* with currency (from) name
 * @param currency_to char* with currency (to) name
 * @return
 */
static bool currency_exists(sqlite3* db, const char* currency_from, const char* currency_to);

/**
 * Returns true if a main_category for given id exists, or false otherwise.
 * @param db sqlite3* database pointer
 * @param maincategory_id main_category_id unsigned long with id of the account.
 * @return true or false
 */
static bool maincategory_exists(sqlite3* db, const unsigned long maincategory_id);

/**
 * Returns true if a transaction for given id exists, or false otherwise.
 * @param db sqlite3* database pointer
 * @param transationc_id unsigned long with id of the account.
 * @return true or false
 */
static bool transaction_exists(sqlite3* db, const unsigned long transaction_id);

/* Functions definitions */

int account_add(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_add_account[SQL_SIZE] = {NULL_STRING};
    char *zErrMsg;
    int result = 0;
    ACCOUNT_TYPE account_type;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name
                , OPTION_FILE_SHORT
                , OPTION_FILE_LONG);
        return 1;
    } else if (parameters.name[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name
                , OPTION_FILE_SHORT
                , OPTION_FILE_LONG);
        return 1;
    } else if (parameters.currency[0] == NULL_STRING) {
        if (parameters.default_currency[0] != NULL_STRING) {
            strncpy(parameters.currency, parameters.default_currency, PAR_CURRENCY_LEN - 1);
        } else {
            fprintf(stderr, MSG_MISSING_PAR_CURRENCY
                    , parameters.prog_name
                    , OPTION_CURRENCY_SHORT
                    , OPTION_CURRENCY_LONG);
            return 1;
        }
    }

    // Set default values
    account_type = (parameters.account_type == ACC_TYPE_UNSET) ? ACC_TRANSACTIONAL : parameters.account_type;

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Prepare SQL query and perform database actions
    sprintf(sql_add_account,
            "INSERT INTO ACCOUNTS VALUES ("
            "null, " // id - auto increment
            "'%s', " // 1 - name
            "'%s', " // 2 - description
            "'%s', " // 3 - institution
            "%d, " // 4 - account type
            "'%s', " // 5 - currency
            "%d);" // 6 - account status
            , parameters.name
            , parameters.description
            , parameters.institution
            , account_type
            , parameters.currency
            , ITEM_STAT_OPEN);

    if (sqlite3_exec(db, sql_add_account, NULL, NULL, &zErrMsg) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, zErrMsg);
        sqlite3_free(zErrMsg);
        result = 1;
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // if succeded and verbose then inform about performed actions
    if (parameters.verbose) {
        if (result == 0) {
            printf("%s: added new account: %s\n"
                    , parameters.prog_name
                    , parameters.name);
        } else {
            printf("%s: no accounts added.\n", parameters.prog_name);
        }
    }

    return result;
}

int account_edit(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_edit_account[SQL_SIZE] = {NULL_STRING};
    char buf[BUF_SIZE];
    char *zErrMsg;
    int result = 0;
    char final_message[BUF_SIZE] = {NULL_STRING};


    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name
                , OPTION_FILE_SHORT
                , OPTION_FILE_LONG);
        return 1;
    } else if (parameters.id == PAR_ID_NOT_SET) {
        fprintf(stderr, MSG_MISSING_PAR_ID
                , parameters.prog_name
                , OPTION_ID_SHORT
                , OPTION_ID_LONG);
        return 1;
    } else if (parameters.currency[0] != NULL_STRING) {
        fprintf(stderr, MSG_CURRENCY_CANNOT_CHANGE
                , parameters.prog_name
                , OPTION_CURRENCY_SHORT
                , OPTION_CURRENCY_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Check if the account exists
    if (account_exists(db, parameters.id) == true) {

        // Prepare SQL query and perform database actions
        strcat(sql_edit_account, "BEGIN TRANSACTION;");

        if (parameters.name[0] != NULL_STRING) {
            sprintf(buf, "UPDATE ACCOUNTS SET NAME='%s' WHERE ACCOUNT_ID=%li;"
                    , parameters.name
                    , parameters.id);
            strncat(sql_edit_account, buf, BUF_SIZE);
        }

        if (parameters.description[0] != NULL_STRING) {
            sprintf(buf, "UPDATE ACCOUNTS SET DESCRIPTION='%s' WHERE ACCOUNT_ID=%li;"
                    , parameters.description
                    , parameters.id);
            strncat(sql_edit_account, buf, BUF_SIZE);
        }

        if (parameters.institution[0] != NULL_STRING) {
            sprintf(buf, "UPDATE ACCOUNTS SET INSTITUTION='%s' WHERE ACCOUNT_ID=%li;"
                    , parameters.institution
                    , parameters.id);
            strncat(sql_edit_account, buf, BUF_SIZE);
        }

        if (parameters.account_type != ACC_TYPE_UNSET) {
            sprintf(buf, "UPDATE ACCOUNTS SET TYPE=%d WHERE ACCOUNT_ID=%li;"
                    , parameters.account_type
                    , parameters.id);
            strncat(sql_edit_account, buf, BUF_SIZE);
        }

        strcat(sql_edit_account, "COMMIT;");

        if (sqlite3_exec(db, sql_edit_account, NULL, NULL, &zErrMsg) != SQLITE_OK) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, zErrMsg);
            sqlite3_free(zErrMsg);
            result = 1;
            sprintf(final_message, "%s: no change done for account with id=%li\n"
                    , parameters.prog_name
                    , parameters.id);
        } else {
            sprintf(final_message, "%s: edited account with id=%li\n"
                    , parameters.prog_name
                    , parameters.id);
        }
    } else {
        result = 1;
        sprintf(final_message, "%s: there is no account with id=%li - no change done\n"
                , parameters.prog_name
                , parameters.id);
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
        sprintf(final_message, "%s: edited account with id=%li\n"
                , parameters.prog_name
                , parameters.id);
    }

    // Inform about performed actions if succeded and verbose
    if (parameters.verbose) {
        printf(final_message);
    }

    return result;

}

int account_close(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_remove_account[SQL_SIZE] = {NULL_STRING};
    char *zErrMsg;
    int result = 0;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name
                , OPTION_FILE_SHORT
                , OPTION_FILE_LONG);
        return 1;
    } else if (parameters.id == PAR_ID_NOT_SET) {
        fprintf(stderr, MSG_MISSING_PAR_ID
                , parameters.prog_name
                , OPTION_ID_SHORT
                , OPTION_ID_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Prepare SQL query and perform database actions
    sprintf(sql_remove_account, "UPDATE ACCOUNTS SET STATUS=%d WHERE ACCOUNT_ID=%li;"
            , ITEM_STAT_CLOSED
            , parameters.id);

    if (sqlite3_exec(db, sql_remove_account, NULL, NULL, &zErrMsg) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, zErrMsg);
        sqlite3_free(zErrMsg);
        result = 1;
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // if succeded and verbose then inform about performed actions
    if (parameters.verbose) {
        if (result == 0) {
            printf("%s: closed account: %li\n", parameters.prog_name, parameters.id);
        } else {

            printf("%s: no accounts closed.\n", parameters.prog_name);
        }
    }

    return result;
}

int account_list(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_list_accounts[SQL_SIZE] = {NULL_STRING};
    char sql_buf[BUF_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    int result = 0;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name
                , OPTION_FILE_SHORT
                , OPTION_FILE_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Prepare SQL query and perform database actions
    sprintf(sql_list_accounts, "SELECT"
            " ACCOUNT_ID, NAME, TYPE, CURRENCY, INSTITUTION, DESCRIPTION"
            " FROM ACCOUNTS WHERE STATUS=%d"
            , ITEM_STAT_OPEN);
    if (parameters.account_type != ACC_TYPE_UNSET) {
        sprintf(sql_buf, " AND TYPE=%d", parameters.account_type);
        strncat(sql_list_accounts, sql_buf, SQL_SIZE);
    }
    if (parameters.currency[0] != NULL_STRING) {
        sprintf(sql_buf, " AND CURRENCY='%s'", parameters.currency);
        strncat(sql_list_accounts, sql_buf, SQL_SIZE);
    }
    if (parameters.institution[0] != NULL_STRING) {
        sprintf(sql_buf, " AND INSTITUTION LIKE '%%%s%%'"
                , parameters.institution);
        strncat(sql_list_accounts, sql_buf, SQL_SIZE);
    }
    strncat(sql_list_accounts, " ORDER BY NAME ASC;", BUF_SIZE);

    if (sqlite3_prepare_v2(db, sql_list_accounts, SQL_SIZE, &sqlStmt, NULL) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        return 1;
    }

    printf(FS_ID_T FS_GAP FS_NAME_T FS_GAP FS_ATYPE_T FS_GAP FS_CUR_T FS_GAP FS_INST_T FS_GAP FS_DESC_T "\n"
            , "ID", "NAME", "TYPE", "CUR", "BANK", "DESCRIPTION");

    // Print accounts on standard output
    while ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
        printf(FS_ID, sqlite3_column_int(sqlStmt, 0));
        printf(FS_GAP FS_NAME, sqlite3_column_text(sqlStmt, 1));
        printf(FS_GAP FS_ATYPE, account_type_text(sqlite3_column_int(sqlStmt, 2)));
        printf(FS_GAP FS_CUR, sqlite3_column_text(sqlStmt, 3));
        printf(FS_GAP FS_INST, sqlite3_column_text(sqlStmt, 4));
        printf(FS_GAP FS_DESC, sqlite3_column_text(sqlStmt, 5));
        printf("\n");
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

int category_add(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_add_category[SQL_SIZE] = {NULL_STRING};
    char *zErrMsg;
    int result = 0;
    long maincategory_id;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name
                , OPTION_FILE_SHORT
                , OPTION_FILE_LONG);
        return 1;
    } else if (parameters.name[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name
                , OPTION_FILE_SHORT
                , OPTION_FILE_LONG);
        return 1;
    } else if (parameters.maincat_name[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_MAINCATEGORY
                , parameters.prog_name
                , OPTION_MAINCATEGORY_SHORT
                , OPTION_MAINCATEGORY_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    maincategory_id = maincategory_id_for_name(db, parameters.maincat_name);
    if (maincategory_id == -1) {
        fprintf(stderr, MSG_MAINCATEGORY_NOT_FOUND
                , parameters.prog_name
                , parameters.maincat_name);
        sqlite3_close(db);
        return 1;
    }

    // Prepare SQL query and perform database actions
    sprintf(sql_add_category,
            "INSERT INTO CATEGORIES VALUES ("
            "null, " // id - auto increment
            "%li, " // 1 - main category
            "'%s', " // 2 - name
            "%d);", // 3 - status
            maincategory_id,
            parameters.name,
            ITEM_STAT_OPEN);

    if (sqlite3_exec(db, sql_add_category, NULL, NULL, &zErrMsg) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, zErrMsg);
        sqlite3_free(zErrMsg);
        result = 1;
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // if succeded and verbose then inform about performed actions
    if (parameters.verbose) {
        if (result == 0) {
            printf("%s: added new category: %s\n"
                    , parameters.prog_name, parameters.name);
        } else {
            printf("%s: no category added.\n", parameters.prog_name);
        }
    }

    return result;
}

int category_edit(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_edit_category[SQL_SIZE] = {NULL_STRING};
    char buf[BUF_SIZE];
    char *zErrMsg;
    int result = 0;
    char final_message[BUF_SIZE] = {NULL_STRING};
    long maincategory_id;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name
                , OPTION_FILE_SHORT
                , OPTION_FILE_LONG);
        return 1;
    } else if (parameters.id == PAR_ID_NOT_SET) {
        fprintf(stderr, MSG_MISSING_PAR_ID
                , parameters.prog_name
                , OPTION_ID_SHORT
                , OPTION_ID_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Check if the category exists
    if (category_exists(db, parameters.id) == true) {

        // Prepare SQL query and perform database actions
        strcat(sql_edit_category, "BEGIN TRANSACTION;");
        if (parameters.maincat_name[0] != NULL_STRING) {
            maincategory_id = maincategory_id_for_name(db, parameters.maincat_name);
            if (maincategory_id == -1) {
                fprintf(stderr, MSG_MAINCATEGORY_NOT_FOUND
                        , parameters.prog_name
                        , parameters.maincat_name);
                sqlite3_close(db);
                return 1;
            } else {
                sprintf(buf, "UPDATE CATEGORIES SET MAIN_CATEGORY_ID=%li"
                        " WHERE CATEGORY_ID=%li;",
                        maincategory_id,
                        parameters.id);
                strncat(sql_edit_category, buf, BUF_SIZE);
            }
        }
        if (parameters.name[0] != NULL_STRING) {
            sprintf(buf, "UPDATE CATEGORIES SET NAME='%s' WHERE CATEGORY_ID=%li;",
                    parameters.name,
                    parameters.id);
            strncat(sql_edit_category, buf, BUF_SIZE);
        }
        strcat(sql_edit_category, "COMMIT;");

        if (sqlite3_exec(db, sql_edit_category, NULL, NULL, &zErrMsg) != SQLITE_OK) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, zErrMsg);
            sqlite3_free(zErrMsg);
            result = 1;
            sprintf(final_message, "%s: no change done for category with id=%li\n"
                    , parameters.prog_name
                    , parameters.id);
        } else {
            sprintf(final_message, "%s: edited category with id=%li\n"
                    , parameters.prog_name
                    , parameters.id);
        }
    } else {
        result = 1;
        sprintf(final_message, "%s: there is no category with id=%li - no change done\n"
                , parameters.prog_name
                , parameters.id);
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
        sprintf(final_message, "%s: edited category with id =%li\n"
                , parameters.prog_name
                , parameters.id);
    }

    // Inform about performed actions if succeeded and verbose
    if (parameters.verbose) {
        printf(final_message);
    }

    return result;
}

int category_remove(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_remove_category[SQL_SIZE] = {NULL_STRING};
    char *zErrMsg;
    int result = 0;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.id == PAR_ID_NOT_SET) {
        fprintf(stderr, MSG_MISSING_PAR_ID
                , parameters.prog_name, OPTION_ID_SHORT, OPTION_ID_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Prepare SQL query and perform database actions
    sprintf(sql_remove_category, "UPDATE CATEGORIES SET STATUS=%d"
            " WHERE CATEGORY_ID=%li;"
            , ITEM_STAT_CLOSED
            , parameters.id);

    if (sqlite3_exec(db, sql_remove_category, NULL, NULL, &zErrMsg) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, zErrMsg);
        sqlite3_free(zErrMsg);
        result = 1;
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // if succeded and verbose then inform about performed actions
    if (parameters.verbose) {
        if (result == 0) {
            printf("%s: removed category: %li\n", parameters.prog_name, parameters.id);
        } else {

            printf("%s: no category removed.\n", parameters.prog_name);
        }
    }

    return result;
}

int category_list(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_list_categories[SQL_SIZE] = {NULL_STRING};
    char sql_buf[BUF_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    int result = 0;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        return 1;
    }

    // Prepare SQL query and perform database actions
    sprintf(sql_list_categories, "SELECT"
            " C.CATEGORY_ID" // 0
            ",M.TYPE" // 1
            ",M.NAME" // 2
            ",C.NAME" // 3
            " FROM CATEGORIES C"
            " LEFT JOIN MAIN_CATEGORIES M ON C.MAIN_CATEGORY_ID=M.MAIN_CATEGORY_ID"
            " WHERE C.STATUS=%d", ITEM_STAT_OPEN);
    if (parameters.maincat_name[0] != NULL_STRING) {
        sprintf(sql_buf, " AND M.NAME LIKE '%%%s%%'", parameters.maincat_name);
        strncat(sql_list_categories, sql_buf, BUF_SIZE);
    }
    if (parameters.maincategory_type != CAT_TYPE_NOTSET) {
        sprintf(sql_buf, " AND M.TYPE=%d", parameters.maincategory_type);
        strncat(sql_list_categories, sql_buf, BUF_SIZE);
    }
    strncat(sql_list_categories, " ORDER BY M.TYPE, M.NAME, C.NAME;", BUF_SIZE);

    if (sqlite3_prepare_v2(db, sql_list_categories, SQL_SIZE, &sqlStmt, NULL) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        return 1;
    }

    printf(FS_ID_T FS_GAP FS_CTYPE_T FS_GAP FS_NAME_T FS_GAP FS_NAME_T "\n", "ID", "TYPE", "MAIN", "CATEGORY");

    // Print transactions on standard output
    while ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
        printf(FS_ID, sqlite3_column_int(sqlStmt, 0)); // id
        printf(FS_GAP FS_CTYPE, maincategory_type_text(sqlite3_column_int(sqlStmt, 1))); // type
        printf(FS_GAP FS_NAME, sqlite3_column_text(sqlStmt, 2)); // main category
        printf(FS_GAP FS_NAME, sqlite3_column_text(sqlStmt, 3)); // category
        printf("\n");
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

int currency_add(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_add_currency[SQL_SIZE] = {NULL_STRING};
    char *zErrMsg;
    int result = 0;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.currency[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_CURRENCY
                , parameters.prog_name, OPTION_CURRENCY_SHORT, OPTION_CURRENCY_LONG);
        return 1;
    } else if (parameters.currency_to[0] == NULL_STRING) {
        if (parameters.default_currency[0] != NULL_STRING) {
            strncpy(parameters.currency_to
                    , parameters.default_currency, PAR_CURRENCY_LEN - 1);
        } else {
            fprintf(stderr, MSG_MISSING_PAR_CURRENCY_TO
                    , parameters.prog_name, OPTION_CURRENCY_TO_SHORT, OPTION_CURRENCY_TO_LONG);
            return 1;
        }
    } else if (parameters.value[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_EXCHANGE_RATE
                , parameters.prog_name, OPTION_VALUE_SHORT, OPTION_VALUE_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Prepare SQL query and perform database actions
    sprintf(sql_add_currency,
            "INSERT INTO CURRENCIES VALUES ("
            "'%s'" // currency from
            ",'%s'" // currency to
            ",round(%f,4));", // echange rate
            parameters.currency,
            parameters.currency_to,
            atof(parameters.value));
    if (sqlite3_exec(db, sql_add_currency, NULL, NULL, &zErrMsg) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, zErrMsg);
        sqlite3_free(zErrMsg);
        result = 1;
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // if succeded and verbose then inform about performed actions
    if (parameters.verbose) {
        if (result == 0) {
            printf("%s: added new currency exchange rate %s-%s\n"
                    , parameters.prog_name, parameters.currency, parameters.currency_to);
        } else {
            printf("%s: no currency exchange rate added.\n"
                    , parameters.prog_name);
        }
    }

    return result;
}

int currency_edit(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_edit_currency[SQL_SIZE] = {NULL_STRING};
    char *zErrMsg;
    int result = 0;
    char final_message[BUF_SIZE] = {NULL_STRING};


    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.currency[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_CURRENCY
                , parameters.prog_name, OPTION_CURRENCY_SHORT, OPTION_CURRENCY_LONG);
        return 1;
    } else if (parameters.currency_to[0] == NULL_STRING) {
        if (parameters.default_currency[0] != NULL_STRING) {
            strncpy(parameters.currency_to
                    , parameters.default_currency, PAR_CURRENCY_LEN - 1);
        } else {
            fprintf(stderr, MSG_MISSING_PAR_CURRENCY_TO
                    , parameters.prog_name
                    , OPTION_CURRENCY_TO_SHORT
                    , OPTION_CURRENCY_TO_LONG);
            return 1;
        }
    } else if (parameters.value[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_EXCHANGE_RATE
                , parameters.prog_name, OPTION_VALUE_SHORT, OPTION_VALUE_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Check if the currencies pair exists
    if (currency_exists(db, parameters.currency, parameters.currency_to) == true) {

        // Prepare SQL query and perform database actions
        sprintf(sql_edit_currency,
                "UPDATE CURRENCIES SET EXCHANGE_RATE=round(%f,4)"
                " WHERE"
                " CURRENCY_FROM='%s'"
                " AND CURRENCY_TO='%s';",
                atof(parameters.value),
                parameters.currency,
                parameters.currency_to
                );

        if (sqlite3_exec(db, sql_edit_currency, NULL, NULL, &zErrMsg) != SQLITE_OK) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, zErrMsg);
            sqlite3_free(zErrMsg);
            result = 1;
            sprintf(final_message, "%s: no change done for exchange rate %s-%s\n"
                    , parameters.prog_name
                    , parameters.currency
                    , parameters.currency_to);
        } else {
            sprintf(final_message, "%s: edited exchange rate for %s-%s\n"
                    , parameters.prog_name
                    , parameters.currency
                    , parameters.currency_to);
        }
    } else {
        result = 1;
        sprintf(final_message, "%s: there is no currncies pair: %s-%s - no change done\n"
                , parameters.prog_name
                , parameters.currency
                , parameters.currency_to);
    }
    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
        sprintf(final_message, "%s: edited exchange rate for %s-%s\n"
                , parameters.prog_name
                , parameters.currency
                , parameters.currency_to);
    }

    // Inform about performed actions if succeded and verbose
    if (parameters.verbose) {
        printf(final_message);
    }

    return result;
}

int currency_remove(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_remove_currency[SQL_SIZE] = {NULL_STRING};
    char *zErrMsg;
    int result = 0;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.currency[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_CURRENCY
                , parameters.prog_name
                , OPTION_CURRENCY_SHORT
                , OPTION_CURRENCY_LONG);
        return 1;
    } else if (parameters.currency_to[0] == NULL_STRING) {
        if (parameters.default_currency[0] != NULL_STRING) {
            strncpy(parameters.currency_to
                    , parameters.default_currency, PAR_CURRENCY_LEN - 1);
        } else {
            fprintf(stderr, MSG_MISSING_PAR_CURRENCY_TO
                    , parameters.prog_name
                    , OPTION_CURRENCY_TO_SHORT
                    , OPTION_CURRENCY_TO_LONG);
            return 1;
        }
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Prepare SQL query and perform database actions
    sprintf(sql_remove_currency, "DELETE FROM CURRENCIES"
            " WHERE"
            " CURRENCY_FROM='%s'"
            " AND CURRENCY_TO='%s';",
            parameters.currency,
            parameters.currency_to
            );

    if (sqlite3_exec(db, sql_remove_currency, NULL, NULL, &zErrMsg) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, zErrMsg);
        sqlite3_free(zErrMsg);
        result = 1;
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // if succeded and verbose then inform about performed actions
    if (parameters.verbose) {
        if (result == 0) {
            printf("%s: removed currency exchange rate for: %s-%s\n"
                    , parameters.prog_name
                    , parameters.currency
                    , parameters.currency_to);
        } else {

            printf("%s: no currency exchange rate removed.\n"
                    , parameters.prog_name);
        }
    }

    return result;
}

int currency_list(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_list_currencies[SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    int result = 0;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Prepare SQL query and perform database actions
    sprintf(sql_list_currencies, "SELECT"
            " CURRENCY_FROM, CURRENCY_TO, EXCHANGE_RATE"
            " FROM CURRENCIES ORDER BY CURRENCY_FROM, CURRENCY_TO;");

    if (sqlite3_prepare_v2(db, sql_list_currencies, SQL_SIZE, &sqlStmt, NULL) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        return 1;
    }

    printf(FS_CURL FS_GAP FS_CURL FS_GAP FS_EXCHRT_T "\n", "CURRENCY FROM", "CURRENCY TO", "EXCHANGE RATE");

    // Print accounts on standard output
    while ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
        printf(FS_CURL, sqlite3_column_text(sqlStmt, 0));
        printf(FS_GAP FS_CURL, sqlite3_column_text(sqlStmt, 1));
        printf(FS_GAP FS_EXCHRT, sqlite3_column_double(sqlStmt, 2));
        printf("\n");
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

int datafile_init(PARAMETERS parameters)
{
    sqlite3 *db;
    char *zErrMsg;
    int result = 0;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Create tables
    char *sqlTransaction = "BEGIN TRANSACTION;"
            "CREATE TABLE CURRENCIES ("
            "CURRENCY_FROM TEXT,"
            "CURRENCY_TO TEXT,"
            "EXCHANGE_RATE REAL,"
            "PRIMARY KEY (CURRENCY_FROM, CURRENCY_TO));"

            "CREATE TABLE ACCOUNTS ("
            "ACCOUNT_ID INTEGER PRIMARY KEY,"
            "NAME TEXT,"
            "DESCRIPTION TEXT,"
            "INSTITUTION TEXT,"
            "TYPE INTEGER,"
            "CURRENCY TEXT,"
            "STATUS INTEGER);"

            "CREATE TABLE TRANSACTIONS ("
            "TRANSACTION_ID INTEGER PRIMARY KEY,"
            "YEAR INTEGER,"
            "MONTH INTEGER,"
            "DAY INTEGER,"
            "ACCOUNT_ID INTEGER,"
            "DESCRIPTION TEXT,"
            "VALUE REAL,"
            "CATEGORY_ID INTEGER);"

            "CREATE TABLE BUDGETS ("
            "YEAR INTEGER,"
            "MONTH INTEGER,"
            "CATEGORY_ID INTEGER,"
            "VALUE REAL,"
            "CURRENCY TEXT,"
            "PRIMARY KEY (YEAR, MONTH, CATEGORY_ID));"

            "CREATE TABLE CATEGORIES ("
            "CATEGORY_ID INTEGER PRIMARY KEY,"
            "MAIN_CATEGORY_ID INTEGER,"
            "NAME TEXT,"
            "STATUS INTEGER);"

            "CREATE TABLE MAIN_CATEGORIES ("
            "MAIN_CATEGORY_ID INTEGER PRIMARY KEY,"
            "TYPE INTEGER,"
            "NAME TEXT,"
            "STATUS INTEGER);"

            "COMMIT;";

    if (sqlite3_exec(db, sqlTransaction, NULL, NULL, &zErrMsg) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, zErrMsg);
        sqlite3_free(zErrMsg);
        result = 1;
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // Inform about performed actions if succeded and verbose
    if (parameters.verbose) {
        if (result == 0) {
            printf("%s: initialized new data file: %s\n"
                    , parameters.prog_name, parameters.dataFilePath);
        } else {

            printf("%s: no file initialized.\n", parameters.prog_name);
        }
    }
    //TODO: function to load currencies exchange rates from internet
    return result;
}

int maincategory_add(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_add_maincategory[SQL_SIZE] = {NULL_STRING};
    char *zErrMsg;
    int result = 0;
    MAINCATEGORY_TYPE maincategory_type;


    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.name[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    }

    // Set default values
    maincategory_type = (parameters.maincategory_type == CAT_TYPE_NOTSET) ? CAT_COST : parameters.maincategory_type;

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }


    // Prepare SQL query and perform database actions
    sprintf(sql_add_maincategory,
            "INSERT INTO MAIN_CATEGORIES VALUES ("
            "null, " // id - auto increment
            "'%d', " // 1 - type
            "'%s', " // 2 - name
            "'%d');", // 3 - status
            maincategory_type,
            parameters.name,
            ITEM_STAT_OPEN);

    if (sqlite3_exec(db, sql_add_maincategory, NULL, NULL, &zErrMsg) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, zErrMsg);
        sqlite3_free(zErrMsg);
        result = 1;
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // if succeded and verbose then inform about performed actions
    if (parameters.verbose) {
        if (result == 0) {
            printf("%s: added new main category: %s\n", parameters.prog_name, parameters.name);
        } else {

            printf("%s: no main categories added.\n", parameters.prog_name);
        }
    }

    return result;
}

int maincategory_edit(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_edit_maincategory[SQL_SIZE] = {NULL_STRING};
    char buf[BUF_SIZE];
    char *zErrMsg;
    int result = 0;
    char final_message[BUF_SIZE] = {NULL_STRING};



    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.id == PAR_ID_NOT_SET) {
        fprintf(stderr, MSG_MISSING_PAR_ID
                , parameters.prog_name, OPTION_ID_SHORT, OPTION_ID_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Check if the maincategory exists
    if (maincategory_exists(db, parameters.id) == true) {

        // Prepare SQL query and perform database actions
        strcat(sql_edit_maincategory, "BEGIN TRANSACTION;");

        if (parameters.maincategory_type != CAT_TYPE_NOTSET) {
            sprintf(buf, "UPDATE MAIN_CATEGORIES SET TYPE=%d WHERE MAIN_CATEGORY_ID=%li;",
                    parameters.maincategory_type,
                    parameters.id);
            strncat(sql_edit_maincategory, buf, BUF_SIZE);
        }
        if (parameters.name[0] != NULL_STRING) {
            sprintf(buf, "UPDATE MAIN_CATEGORIES SET NAME='%s' WHERE MAIN_CATEGORY_ID=%li;",
                    parameters.name,
                    parameters.id);
            strncat(sql_edit_maincategory, buf, BUF_SIZE);
        }

        strcat(sql_edit_maincategory, "COMMIT;");

        if (sqlite3_exec(db, sql_edit_maincategory, NULL, NULL, &zErrMsg) != SQLITE_OK) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, zErrMsg);
            sqlite3_free(zErrMsg);
            result = 1;
            sprintf(final_message, "%s: no change done for main category with id=%li\n"
                    , parameters.prog_name
                    , parameters.id);
        } else {
            sprintf(final_message, "%s: edited main category with id=%li\n"
                    , parameters.prog_name
                    , parameters.id);
        }
    } else {
        result = 1;
        sprintf(final_message, "%s: there is no main category with id=%li - no change done\n"
                , parameters.prog_name
                , parameters.id);
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
        sprintf(final_message, "%s: edited main cateory with id=%li\n"
                , parameters.prog_name
                , parameters.id);
    }

    // Inform about performed actions if succeded and verbose
    if (parameters.verbose) {
        printf(final_message);
    }

    return result;
}

int maincategory_remove(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_remove_maincategory[SQL_SIZE] = {NULL_STRING};
    char *zErrMsg;
    int result = 0;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.id == PAR_ID_NOT_SET) {
        fprintf(stderr, MSG_MISSING_PAR_ID
                , parameters.prog_name, OPTION_ID_SHORT, OPTION_ID_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Prepare SQL query and perform database actions
    sprintf(sql_remove_maincategory, "UPDATE MAIN_CATEGORIES"
            " SET STATUS=%d WHERE MAIN_CATEGORY_ID=%li;"
            , ITEM_STAT_CLOSED
            , parameters.id);

    if (sqlite3_exec(db, sql_remove_maincategory, NULL, NULL, &zErrMsg) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, zErrMsg);
        sqlite3_free(zErrMsg);
        result = 1;
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // if succeded and verbose then inform about performed actions
    if (parameters.verbose) {
        if (result == 0) {
            printf("%s: removed main category: %li\n"
                    , parameters.prog_name, parameters.id);
        } else {

            printf("%s: no main category removed.\n", parameters.prog_name);
        }
    }

    return result;
}

int maincategory_list(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_list_maincategories[SQL_SIZE] = {NULL_STRING};
    char sql_buf[SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    int result = 0;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Prepare SQL query and perform database actions
    sprintf(sql_list_maincategories, "SELECT"
            " MAIN_CATEGORY_ID, TYPE, NAME"
            " FROM MAIN_CATEGORIES WHERE STATUS=%d"
            , ITEM_STAT_OPEN);
    if (parameters.maincategory_type != CAT_TYPE_NOTSET) {
        sprintf(sql_buf, " AND TYPE=%d", parameters.maincategory_type);
        strncat(sql_list_maincategories, sql_buf, BUF_SIZE);
    }
    strncat(sql_list_maincategories, " ORDER BY TYPE, NAME;", BUF_SIZE);

    if (sqlite3_prepare_v2(db, sql_list_maincategories, SQL_SIZE, &sqlStmt, NULL) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        return 1;
    }

    printf(FS_ID_T FS_GAP FS_MTYPE_T FS_GAP FS_NAME_T "\n", "ID", "TYPE", "NAME");

    // Print accounts on standard output
    while ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
        printf(FS_ID, sqlite3_column_int(sqlStmt, 0));
        printf(FS_GAP FS_MTYPE, maincategory_type_text(sqlite3_column_int(sqlStmt, 1)));
        printf(FS_GAP FS_NAME, sqlite3_column_text(sqlStmt, 2));
        printf("\n");
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

int transaction_add(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_add_transaction[SQL_SIZE] = {NULL_STRING};
    char *zErrMsg;
    int result = 0;
    long account_id, category_id;
    int year, month, day;
    float transaction_value;

    // Make sure necessary data have been provided
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.acc_name[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_ACCOUNT_NAME
                , parameters.prog_name, OPTION_ACCOUNT_SHORT, OPTION_ACCOUNT_LONG);
        return 1;
    } else if (parameters.description[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_DESCRIPTION
                , parameters.prog_name, OPTION_DESCRIPTION_SHORT, OPTION_DESCRIPTION_LONG);
        return 1;
    } else if (parameters.value[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_VALUE
                , parameters.prog_name, OPTION_VALUE_SHORT, OPTION_VALUE_LONG);
        return 1;
    } else if (parameters.cat_name[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_CATEGORY
                , parameters.prog_name, OPTION_CATEGORY_SHORT, OPTION_CATEGORY_LONG);
        return 1;
    }


    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }
    // Get neccessary data
    account_id = account_id_for_name(db, parameters.acc_name);
    if (account_id == -1) {
        fprintf(stderr, MSG_ACCOUNT_NOT_FOUND
                , parameters.prog_name, parameters.acc_name);
        sqlite3_close(db);
        return 1;
    }
    category_id = category_id_for_name(db, parameters.cat_name);
    if (category_id == -1) {
        fprintf(stderr, MSG_CATEGORY_NOT_FOUND
                , parameters.prog_name, parameters.cat_name);
        sqlite3_close(db);
        return 1;
    }
    if (parameters.date[0] == NULL_STRING) {
        if (date_from_string(parameters.date_default, &year, &month, &day) != DT_FULL_DATE) {
            fprintf(stderr, MSG_WRONG_DATE_FULL, parameters.prog_name);
            return 1;
        }
    } else {
        if (date_from_string(parameters.date, &year, &month, &day) != DT_FULL_DATE) {
            fprintf(stderr, MSG_WRONG_DATE_FULL, parameters.prog_name);
            return 1;
        }
    }

    transaction_value = category_factor_for_id(db, category_id) * atof(parameters.value);

    // Prepare SQL query and perform database operations
    sprintf(sql_add_transaction,
            "INSERT INTO TRANSACTIONS VALUES ("
            "null, " // id - auto increment
            "%d, " // year
            "%d, " // month
            "%d, " // day
            "%li, " // account id
            "'%s', " // description
            "round(%f,2), " // value
            "%li);", // category id
            year,
            month,
            day,
            account_id,
            parameters.description,
            transaction_value,
            category_id
            );

    if (sqlite3_exec(db, sql_add_transaction, NULL, NULL, &zErrMsg) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, zErrMsg);
        sqlite3_free(zErrMsg);
        result = 1;
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // If succeded and verbose then inform about performed actions
    if (parameters.verbose) {
        if (result == 0) {
            printf("%s: added new transaction.\n", parameters.prog_name);
        } else {

            printf("%s: no transaction added.\n", parameters.prog_name);
        }
    }

    return result;
}

int transaction_edit(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_edit_transaction[SQL_SIZE] = {NULL_STRING};
    char buf[BUF_SIZE];
    char *zErrMsg;
    int result = 0;
    int year, month, day;
    long account_id, category_id;
    float transaction_value;
    char final_message[BUF_SIZE] = {NULL_STRING};


    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.id == PAR_ID_NOT_SET) {
        fprintf(stderr, MSG_MISSING_PAR_ID
                , parameters.prog_name, OPTION_ID_SHORT, OPTION_ID_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Check if the transaction exists
    if (transaction_exists(db, parameters.id) == true) {

        // Prepare SQL query and perform database actions
        strcat(sql_edit_transaction, "BEGIN TRANSACTION;");

        if (parameters.date[0] != NULL_STRING) {
            if (date_from_string(parameters.date, &year, &month, &day) != DT_FULL_DATE) {
                fprintf(stderr, MSG_WRONG_DATE_FULL, parameters.prog_name);
                return 1;
            } else {
                sprintf(buf, "UPDATE TRANSACTIONS SET YEAR=%d WHERE TRANSACTION_ID=%li;"
                        "UPDATE TRANSACTIONS SET MONTH=%d WHERE TRANSACTION_ID=%li;"
                        "UPDATE TRANSACTIONS SET DAY=%d WHERE TRANSACTION_ID=%li;",
                        year, parameters.id,
                        month, parameters.id,
                        day, parameters.id);
                strncat(sql_edit_transaction, buf, BUF_SIZE);
            }
        }

        if (parameters.acc_name[0] != NULL_STRING) {
            account_id = account_id_for_name(db, parameters.acc_name);
            if (account_id == -1) {
                fprintf(stderr, MSG_ACCOUNT_NOT_FOUND
                        , parameters.prog_name, parameters.acc_name);
                sqlite3_close(db);
                return 1;
            } else {
                sprintf(buf, "UPDATE TRANSACTIONS SET ACCOUNT_ID=%li WHERE TRANSACTION_ID=%li;",
                        account_id,
                        parameters.id);
                strncat(sql_edit_transaction, buf, BUF_SIZE);
            }
        }

        if (parameters.description[0] != NULL_STRING) {
            sprintf(buf, "UPDATE TRANSACTIONS SET DESCRIPTION='%s' WHERE TRANSACTION_ID=%li;",
                    parameters.description,
                    parameters.id);
            strncat(sql_edit_transaction, buf, BUF_SIZE);
        }

        if (parameters.cat_name[0] != NULL_STRING) {
            category_id = category_id_for_name(db, parameters.cat_name);
            if (category_id == -1) {
                fprintf(stderr, MSG_CATEGORY_NOT_FOUND
                        , parameters.prog_name, parameters.cat_name);
                sqlite3_close(db);
                return 1;
            } else {
                //TODO: add automatic sign change in case you change category from income to cost or opposite
                sprintf(buf, "UPDATE TRANSACTIONS SET CATEGORY_ID=%li WHERE TRANSACTION_ID=%li;",
                        category_id,
                        parameters.id);
                strncat(sql_edit_transaction, buf, BUF_SIZE);
            }
        } else {
            category_id = transaction_category_id(db, parameters.id);
        }
        if (parameters.value[0] != NULL_STRING) {
            transaction_value = category_factor_for_id(db, category_id) * atof(parameters.value);
            sprintf(buf, "UPDATE TRANSACTIONS SET VALUE=round(%f,2) WHERE TRANSACTION_ID=%li;",
                    transaction_value,
                    parameters.id);
            strncat(sql_edit_transaction, buf, BUF_SIZE);
        }
        strcat(sql_edit_transaction, "COMMIT;");

        if (sqlite3_exec(db, sql_edit_transaction, NULL, NULL, &zErrMsg) != SQLITE_OK) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, zErrMsg);
            sqlite3_free(zErrMsg);
            result = 1;
            sprintf(final_message, "%s: no change done for transaction with id=%li\n"
                    , parameters.prog_name
                    , parameters.id);
        } else {
            sprintf(final_message, "%s: edited transaction with id=%li\n"
                    , parameters.prog_name
                    , parameters.id);
        }
    } else {
        result = 1;
        sprintf(final_message, "%s: there is no transaction with id=%li - no change done\n"
                , parameters.prog_name
                , parameters.id);
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
        sprintf(final_message, "%s: edited transaction with id=%li\n"
                , parameters.prog_name
                , parameters.id);
    }

    // Inform about performed actions if succeeded and verbose
    if (parameters.verbose) {
        printf(final_message);
    }

    return result;
}

int transaction_remove(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_remove_transaction[SQL_SIZE] = {NULL_STRING};
    char *zErrMsg;
    int result = 0;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.id == PAR_ID_NOT_SET) {
        fprintf(stderr, MSG_MISSING_PAR_ID
                , parameters.prog_name, OPTION_ID_SHORT, OPTION_ID_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Prepare SQL query and perform database actions
    sprintf(sql_remove_transaction, "DELETE FROM TRANSACTIONS WHERE TRANSACTION_ID=%li;", parameters.id);
    if (sqlite3_exec(db, sql_remove_transaction, NULL, NULL, &zErrMsg) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, zErrMsg);
        sqlite3_free(zErrMsg);
        result = 1;
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // If succeded and verbose then inform about performed actions
    if (parameters.verbose) {
        if (result == 0) {
            printf("%s: deleted transaction: %li\n", parameters.prog_name, parameters.id);
        } else {

            printf("%s: no transaction deleted.\n", parameters.prog_name);
        }
    }

    return result;
}

int transaction_list(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_list_transactions[SQL_SIZE] = {NULL_STRING};
    char sql_buf[BUF_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    int result = 0;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        return 1;
    }

    // Prepare SQL query and perform database actions
    int where_already_set = 0;
    sprintf(sql_list_transactions, "SELECT"
            " T.TRANSACTION_ID" // 0
            ",T.YEAR" // 1
            ",T.MONTH" // 2
            ",T.DAY" // 3
            ",T.DESCRIPTION" // 4
            ",T.VALUE" // 5
            ",A.NAME" // 6
            ",A.CURRENCY" // 7
            ",C.NAME" // 8
            ",M.NAME" // 9
            ",M.TYPE" // 10
            " FROM TRANSACTIONS T"
            " LEFT JOIN ACCOUNTS A ON T.ACCOUNT_ID=A.ACCOUNT_ID"
            " LEFT JOIN CATEGORIES C ON T.CATEGORY_ID=C.CATEGORY_ID"
            " LEFT JOIN MAIN_CATEGORIES M ON C.MAIN_CATEGORY_ID=M.MAIN_CATEGORY_ID");
    if (parameters.date[0] != NULL_STRING) {
        int year, month, day;
        DATE_TYPE date_type = date_from_string(parameters.date, &year, &month, &day);
        switch (date_type) {
        case DT_FULL_DATE:
            if (where_already_set == 0) {
                sprintf(sql_buf, " WHERE T.YEAR=%d AND T.MONTH=%d AND T.DAY=%d"
                        , year, month, day);
                where_already_set = 1;
            } else {
                sprintf(sql_buf, " AND T.YEAR=%d AND T.MONTH=%d AND T.DAY=%d",
                        year, month, day);
            }
            break;
        case DT_MONTH:
            if (where_already_set == 0) {
                sprintf(sql_buf, " WHERE T.YEAR=%d AND T.MONTH=%d", year, month);
                where_already_set = 1;
            } else {
                sprintf(sql_buf, " AND T.YEAR=%d AND T.MONTH=%d", year, month);
            }
            break;
        case DT_YEAR:
            if (where_already_set == 0) {
                sprintf(sql_buf, " WHERE T.YEAR=%d", year);
                where_already_set = 1;
            } else {
                sprintf(sql_buf, "AND T.YEAR=%d", year);
            }
            break;
        default:
            fprintf(stderr, MSG_WRONG_DATE_FULL, parameters.prog_name);
            return 1;
        }
        strncat(sql_list_transactions, sql_buf, BUF_SIZE);
    }
    if (parameters.acc_name[0] != NULL_STRING) {
        if (where_already_set == 0) {
            sprintf(sql_buf, " WHERE A.NAME LIKE '%%%s%%'", parameters.acc_name);
            where_already_set = 1;
        } else {
            sprintf(sql_buf, " AND A.NAME LIKE '%%%s%%'", parameters.acc_name);
        }
        strncat(sql_list_transactions, sql_buf, BUF_SIZE);
    }
    if (parameters.cat_name[0] != NULL_STRING) {
        if (where_already_set == 0) {
            sprintf(sql_buf, " WHERE C.NAME LIKE '%%%s%%'", parameters.cat_name);
            where_already_set = 1;
        } else {
            sprintf(sql_buf, " AND C.NAME LIKE '%%%s%%'", parameters.cat_name);
        }
        strncat(sql_list_transactions, sql_buf, BUF_SIZE);
    }
    if (parameters.maincat_name[0] != NULL_STRING) {
        if (where_already_set == 0) {
            sprintf(sql_buf, " WHERE M.NAME LIKE '%%%s%%'", parameters.maincat_name);
            where_already_set = 1;
        } else {
            sprintf(sql_buf, " AND M.NAME LIKE '%%%s%%'", parameters.maincat_name);
        }
        strncat(sql_list_transactions, sql_buf, BUF_SIZE);
    }
    if (parameters.maincategory_type != CAT_TYPE_NOTSET) {
        if (where_already_set == 0) {
            sprintf(sql_buf, " WHERE M.TYPE=%d", parameters.maincategory_type);
            where_already_set = 1;
        } else {
            sprintf(sql_buf, " AND M.TYPE=%d", parameters.maincategory_type);
        }
        strncat(sql_list_transactions, sql_buf, BUF_SIZE);
    }
    if (parameters.currency[0] != NULL_STRING) {
        if (where_already_set == 0) {
            sprintf(sql_buf, " WHERE A.CURRENCY='%s'", parameters.currency);
            where_already_set = 1;
        } else {
            sprintf(sql_buf, " AND A.CURRENCY='%s'", parameters.currency);
        }
        strncat(sql_list_transactions, sql_buf, BUF_SIZE);
    }

    strncat(sql_list_transactions, " ORDER BY T.YEAR, T.MONTH, T.DAY;", BUF_SIZE);

    if (sqlite3_prepare_v2(db, sql_list_transactions, SQL_SIZE, &sqlStmt, NULL) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        return 1;
    }

    printf(FS_ID_T FS_GAP FS_DATE_T FS_GAP FS_NAME_T FS_GAP FS_MTYPE_T FS_GAP FS_NAME_T FS_GAP FS_NAME_T FS_GAP FS_VALUE_T FS_GAPS FS_CUR_T FS_GAP FS_DESC_T "\n","ID", "DATE", "ACCOUNT", "TYPE", "MAIN CAT.", "CATEGORY", "VALUE", "CUR", "DESCRIPTION");

    // Print transactions on standard output
    while ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
        printf(FS_ID, sqlite3_column_int(sqlStmt, 0)); // id
        printf(FS_GAP FS_DATE, sqlite3_column_int(sqlStmt, 1), sqlite3_column_int(sqlStmt, 2), sqlite3_column_int(sqlStmt, 3)); // date
        printf(FS_GAP FS_NAME, sqlite3_column_text(sqlStmt, 6)); // account
        printf(FS_GAP FS_MTYPE, maincategory_type_text(sqlite3_column_int(sqlStmt, 10)));
        printf(FS_GAP FS_NAME, sqlite3_column_text(sqlStmt, 9)); // main category
        printf(FS_GAP FS_NAME, sqlite3_column_text(sqlStmt, 8)); // category
        printf(FS_GAP FS_VALUE, sqlite3_column_double(sqlStmt, 5)); // value
        printf(FS_GAPS FS_CUR, sqlite3_column_text(sqlStmt, 7)); // currency
        printf(FS_GAP FS_DESC, sqlite3_column_text(sqlStmt, 4)); // description
        printf("\n");
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

int budget_add(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_add_budget[SQL_SIZE] = {NULL_STRING};
    char *zErrMsg;
    int result = 0;
    long category_id;
    int year, month;
    float budget_value;

    // Make sure necessary data have been provided
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.date[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_DATE
                , parameters.prog_name, OPTION_DATE_SHORT, OPTION_DATE_LONG);
        return 1;
    } else if (parameters.cat_name[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_CATEGORY
                , parameters.prog_name, OPTION_CATEGORY_SHORT, OPTION_CATEGORY_LONG);
        return 1;
    } else if (parameters.value[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_VALUE
                , parameters.prog_name, OPTION_VALUE_SHORT, OPTION_VALUE_LONG);
        return 1;
    } else if (parameters.currency[0] == NULL_STRING) {
        if (parameters.default_currency[0] != NULL_STRING) {
            strncpy(parameters.currency
                    , parameters.default_currency, PAR_CURRENCY_LEN - 1);
        } else {
            fprintf(stderr, MSG_MISSING_PAR_CURRENCY
                    , parameters.prog_name, OPTION_CURRENCY_SHORT, OPTION_CURRENCY_LONG);
            return 1;
        }
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Get necessary data
    category_id = category_id_for_name(db, parameters.cat_name);
    if (category_id == -1) {
        fprintf(stderr, MSG_CATEGORY_NOT_FOUND
                , parameters.prog_name, parameters.cat_name);
        sqlite3_close(db);
        return 1;
    }
    if (date_from_string(parameters.date, &year, &month, NULL) != DT_MONTH) {
        fprintf(stderr, MSG_WRONG_DATE_MONTH, parameters.prog_name);
        return 1;
    }
    budget_value = category_factor_for_id(db, category_id) * atof(parameters.value);

    // Prepare SQL query and perform database operations
    sprintf(sql_add_budget,
            "INSERT INTO BUDGETS VALUES ("
            "%d, " // year
            "%d, " // month
            "%li, " // category id
            "round(%f,2), " // value
            "'%s'" // currency
            ");",
            year,
            month,
            category_id,
            budget_value,
            parameters.currency
            );

    if (sqlite3_exec(db, sql_add_budget, NULL, NULL, &zErrMsg) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, zErrMsg);
        sqlite3_free(zErrMsg);
        result = 1;
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // If succeded and verbose then inform about performed actions
    if (parameters.verbose) {
        if (result == 0) {
            printf("%s: added new budget.\n", parameters.prog_name);
        } else {

            printf("%s: no budget added.\n", parameters.prog_name);
        }
    }

    return result;
}

int budget_list(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_list_budgets[SQL_SIZE] = {NULL_STRING};
    char sql_buf[BUF_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    int result = 0;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        return 1;
    }

    // Prepare SQL query and perform database actions
    int where_already_set = 0;
    sprintf(sql_list_budgets, "SELECT"
            " b.YEAR" // 0
            ", b.MONTH" // 1
            ", mc.TYPE" // 2
            ", mc.NAME" // 3
            ", c.NAME" // 4
            ", b.VALUE" // 5
            ", b.CURRENCY" //6
            " FROM BUDGETS b"
            " LEFT JOIN CATEGORIES c ON b.CATEGORY_ID=c.CATEGORY_ID"
            " LEFT JOIN MAIN_CATEGORIES mc ON c.MAIN_CATEGORY_ID=mc.MAIN_CATEGORY_ID");
    if (parameters.cat_name[0] != NULL_STRING) {
        if (where_already_set == 0) {
            sprintf(sql_buf, " WHERE c.NAME LIKE '%%%s%%'", parameters.cat_name);
            where_already_set = 1;
        } else {
            sprintf(sql_buf, " AND c.NAME LIKE '%%%s%%'", parameters.cat_name);
        }
        strncat(sql_list_budgets, sql_buf, BUF_SIZE);
    }
    if (parameters.date[0] != NULL_STRING) {
        int year, month;
        DATE_TYPE date_type = date_from_string(parameters.date, &year, &month, NULL);
        switch (date_type) {
        case DT_YEAR:
             if (where_already_set == 0) {
                  sprintf(sql_buf, " WHERE b.YEAR=%d", year);
                  where_already_set = 1;
             } else {
                  sprintf(sql_buf, " AND b.YEAR=%d", year);
             }
             break;
        case DT_MONTH:
             if (where_already_set == 0) {
                  sprintf(sql_buf, " WHERE b.YEAR=%d AND b.MONTH=%d", year, month);
                  where_already_set = 1;
             } else {
                  sprintf(sql_buf, " AND b.YEAR=%d AND b.MONTH=%d", year, month);
             }
             break;
        default:
             fprintf(stderr, MSG_WRONG_DATE_MONTH_OR_YEAR, parameters.prog_name);
             return 1;
        }
        strncat(sql_list_budgets, sql_buf, BUF_SIZE);
    }
    strncat(sql_list_budgets, " ORDER BY b.YEAR, b.MONTH, mc.TYPE, mc.NAME, c.NAME;", BUF_SIZE);

    if (sqlite3_prepare_v2(db, sql_list_budgets, SQL_SIZE, &sqlStmt, NULL) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        return 1;
    }

    printf(FS_MONTH_T FS_GAP FS_MTYPE_T FS_GAP FS_NAME_T FS_GAP FS_NAME_T FS_GAP FS_VALUE_T FS_GAPS FS_CUR_T "\n",
            "MONTH", "TYPE", "MAIN CAT.", "CATEGORY", "LIMIT", "CUR");

    // Print transactions on standard output
    while ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
        printf(FS_MONTH, sqlite3_column_int(sqlStmt, 0), sqlite3_column_int(sqlStmt, 1)); // month
        printf(FS_GAP FS_MTYPE, maincategory_type_text(sqlite3_column_int(sqlStmt, 2)));
        printf(FS_GAP FS_NAME, sqlite3_column_text(sqlStmt, 3)); // main category
        printf(FS_GAP FS_NAME, sqlite3_column_text(sqlStmt, 4)); // category
        printf(FS_GAP FS_VALUE, sqlite3_column_double(sqlStmt, 5)); // value
        printf(FS_GAPS FS_CUR, sqlite3_column_text(sqlStmt, 6)); // currency
        printf("\n");
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

int budget_edit(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_edit_budget[SQL_SIZE] = {NULL_STRING};
    char buf[BUF_SIZE];
    char *zErrMsg;
    int result = 0;
    int year, month;
    long category_id;
    float budget_value;
    char final_message[BUF_SIZE] = {NULL_STRING};

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.date[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_DATE
                , parameters.prog_name, OPTION_DATE_SHORT, OPTION_DATE_LONG);
        return 1;
    } else if (parameters.cat_name[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_CATEGORY
                , parameters.prog_name, OPTION_CATEGORY_SHORT, OPTION_CATEGORY_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Get necessary data
    category_id = category_id_for_name(db, parameters.cat_name);
    if (category_id == -1) {
        fprintf(stderr, MSG_CATEGORY_NOT_FOUND
                , parameters.prog_name, parameters.cat_name);
        sqlite3_close(db);
        return 1;
    }
    if (date_from_string(parameters.date, &year, &month, NULL) != DT_MONTH) {
        fprintf(stderr, MSG_WRONG_DATE_MONTH, parameters.prog_name);
        return 1;
    }

    // Check if the budget exists
    if (budget_exists(db, category_id, year, month) == true) {

        // Prepare SQL query and perform database actions
        strcat(sql_edit_budget, "BEGIN TRANSACTION;");

        if (parameters.value[0] != NULL_STRING) {
            budget_value = category_factor_for_id(db, category_id) * atof(parameters.value);
            sprintf(buf, "UPDATE BUDGETS SET VALUE=round(%f,2)"
                    " WHERE YEAR=%d AND MONTH=%d AND CATEGORY_ID=%li;",
                    budget_value,
                    year, month, category_id);
            strncat(sql_edit_budget, buf, BUF_SIZE);
        }
        if (parameters.currency[0] != NULL_STRING) {
            sprintf(buf, "UPDATE BUDGETS SET CURRENCY='%s'"
                    " WHERE YEAR=%d AND MONTH=%d AND CATEGORY_ID=%li;",
                    parameters.currency,
                    year, month, category_id);
            strncat(sql_edit_budget, buf, BUF_SIZE);
        }

        strcat(sql_edit_budget, "COMMIT;");

        if (sqlite3_exec(db, sql_edit_budget, NULL, NULL, &zErrMsg) != SQLITE_OK) {
            fprintf(stderr, "%s: %s\n", parameters.prog_name, zErrMsg);
            sqlite3_free(zErrMsg);
            result = 1;
            sprintf(final_message, "%s: no change done for budget %d-%d %s\n"
                    , parameters.prog_name, year, month, parameters.cat_name);
        } else {
            sprintf(final_message, "%s: edited budget %d-%d %s\n"
                    , parameters.prog_name, year, month, parameters.cat_name);
        }
    } else {
        result = 1;
        sprintf(final_message, "%s: there is no budget %d-%d %s - no change done\n"
                , parameters.prog_name, year, month, parameters.cat_name);
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
        sprintf(final_message, "%s: edited budget %d-%d %s\n"
                , parameters.prog_name, year, month, parameters.cat_name);
    }

    // Inform about performed actions if succeeded and verbose
    if (parameters.verbose) {
        printf(final_message);
    }

    return result;
}

int budget_remove(PARAMETERS parameters)
{
    sqlite3 *db;
    char sql_remove_budget[SQL_SIZE] = {NULL_STRING};
    char *zErrMsg;
    long category_id;
    int year, month;
    int result = 0;

    // Make sure necessary data have been delivered
    if (parameters.dataFilePath[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_FILE
                , parameters.prog_name, OPTION_FILE_SHORT, OPTION_FILE_LONG);
        return 1;
    } else if (parameters.date[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_DATE
                , parameters.prog_name, OPTION_DATE_SHORT, OPTION_DATE_LONG);
        return 1;
    } else if (parameters.cat_name[0] == NULL_STRING) {
        fprintf(stderr, MSG_MISSING_PAR_CATEGORY
                , parameters.prog_name, OPTION_CATEGORY_SHORT, OPTION_CATEGORY_LONG);
        return 1;
    }

    // Open database file
    if (sqlite3_open(parameters.dataFilePath, &db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        sqlite3_close(db);
        return 1;
    }

    // Get necessary data
    category_id = category_id_for_name(db, parameters.cat_name);
    if (category_id == -1) {
        fprintf(stderr, MSG_CATEGORY_NOT_FOUND, parameters.prog_name, parameters.cat_name);
        sqlite3_close(db);
        return 1;
    }
    if (date_from_string(parameters.date, &year, &month, NULL) != DT_MONTH) {
        fprintf(stderr, MSG_WRONG_DATE_MONTH, parameters.prog_name);
        return 1;
    }

    // Prepare SQL query and perform database actions
    sprintf(sql_remove_budget, "DELETE FROM BUDGETS WHERE"
            " YEAR=%d"
            " AND MONTH=%d"
            " AND CATEGORY_ID=%li"
            ";"
            , year, month, category_id);
    if (sqlite3_exec(db, sql_remove_budget, NULL, NULL, &zErrMsg) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, zErrMsg);
        sqlite3_free(zErrMsg);
        result = 1;
    }

    // Close database file
    if (sqlite3_close(db) != SQLITE_OK) {
        fprintf(stderr, "%s: %s\n", parameters.prog_name, sqlite3_errmsg(db));
        result = 1;
    }

    // If succeded and verbose then inform about performed actions
    if (parameters.verbose) {
        if (result == 0) {
            printf("%s: deleted budget.\n", parameters.prog_name);
        } else {

            printf("%s: no budget deleted.\n", parameters.prog_name);
        }
    }
    //TODO: make possible to remove budget for a month/year
    return result;
}

/* Local functions */

static bool account_exists(sqlite3* db, const unsigned long account_id)
{
    char sql_account_exist[SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    bool result = false;

    sprintf(sql_account_exist
            , "SELECT count(ACCOUNT_ID) FROM ACCOUNTS WHERE ACCOUNT_ID=%li AND STATUS=%d;"
            , account_id
            , ITEM_STAT_OPEN
            );

    if (sqlite3_prepare_v2(db, sql_account_exist, SQL_SIZE, &sqlStmt, NULL) == SQLITE_OK) {
        if ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
            if (sqlite3_column_int(sqlStmt, 0) != 0) {
                result = true;
            }
        }
    }

    // Clean
    rc = sqlite3_finalize(sqlStmt);

    return result;

}

static bool budget_exists(sqlite3* db, const unsigned long category_id, const int year, const int month)
{
    char sql_budget_exists[SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    bool result = false;

    sprintf(sql_budget_exists
            , "SELECT count(YEAR) FROM BUDGETS WHERE YEAR=%d AND MONTH=%d AND CATEGORY_ID=%li;"
            , year, month, category_id);

    if (sqlite3_prepare_v2(db, sql_budget_exists, SQL_SIZE, &sqlStmt, NULL) == SQLITE_OK) {
        if ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
            if (sqlite3_column_int(sqlStmt, 0) != 0) {
                result = true;
            }
        }
    }

    // Clean
    rc = sqlite3_finalize(sqlStmt);

    return result;
}

static bool category_exists(sqlite3* db, const unsigned long category_id)
{
    char sql_category_exists[SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    bool result = false;

    sprintf(sql_category_exists
            , "SELECT count(CATEGORY_ID) FROM CATEGORIES WHERE CATEGORY_ID=%li AND STATUS=%d;"
            , category_id
            , ITEM_STAT_OPEN
            );

    if (sqlite3_prepare_v2(db, sql_category_exists, SQL_SIZE, &sqlStmt, NULL) == SQLITE_OK) {
        if ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
            if (sqlite3_column_int(sqlStmt, 0) != 0) {
                result = true;
            }
        }
    }

    // Clean
    rc = sqlite3_finalize(sqlStmt);

    return result;

}

static bool currency_exists(sqlite3* db, const char* currency_from, const char* currency_to)
{
    char sql_currency_exists[SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    bool result = false;

    sprintf(sql_currency_exists
            , "SELECT count(CURRENCY_FROM) FROM CURRENCIES"
            " WHERE CURRENCY_FROM='%s' AND CURRENCY_TO='%s';"
            , currency_from
            , currency_to
            );

    if (sqlite3_prepare_v2(db, sql_currency_exists, SQL_SIZE, &sqlStmt, NULL) == SQLITE_OK) {
        if ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
            if (sqlite3_column_int(sqlStmt, 0) != 0) {
                result = true;
            }
        }
    }

    // Clean
    rc = sqlite3_finalize(sqlStmt);

    return result;
}

static bool maincategory_exists(sqlite3* db, const unsigned long maincategory_id)
{
    char sql_maincategory_exists[SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    bool result = false;

    sprintf(sql_maincategory_exists
            , "SELECT count(MAIN_CATEGORY_ID) FROM MAIN_CATEGORIES WHERE MAIN_CATEGORY_ID=%li AND STATUS=%d;"
            , maincategory_id
            , ITEM_STAT_OPEN
            );

    if (sqlite3_prepare_v2(db, sql_maincategory_exists, SQL_SIZE, &sqlStmt, NULL) == SQLITE_OK) {
        if ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
            if (sqlite3_column_int(sqlStmt, 0) != 0) {
                result = true;
            }
        }
    }

    // Clean
    rc = sqlite3_finalize(sqlStmt);

    return result;

}

static bool transaction_exists(sqlite3* db, const unsigned long transaction_id)
{
    char sql_transaction_exists[SQL_SIZE] = {NULL_STRING};
    sqlite3_stmt *sqlStmt;
    int rc;
    bool result = false;

    sprintf(sql_transaction_exists
            , "SELECT count(TRANSACTION_ID) FROM TRANSACTIONS WHERE TRANSACTION_ID=%li;"
            , transaction_id
            );

    if (sqlite3_prepare_v2(db, sql_transaction_exists, SQL_SIZE, &sqlStmt, NULL) == SQLITE_OK) {
        if ((rc = sqlite3_step(sqlStmt)) == SQLITE_ROW) {
            if (sqlite3_column_int(sqlStmt, 0) != 0) {
                result = true;
            }
        }
    }

    // Clean
    rc = sqlite3_finalize(sqlStmt);

    return result;

}


//TODO: REVIEW: make 'read-only' parameters const
