/*
  Written 2015 by Marcin 'Zbroju' Zbroinski.
  Use of this source code is governed by a GNU General Public License
  that can be found in the LICENSE file.
*/

/* INCLUDE GUARD */
#ifndef MESSAGES_H
#define MESSAGES_H

#include <sqlite3.h>
#include <stdbool.h>

/* MESSAGES */
#define MSG_MISSING_PAR_FILE "%s: missing data file path. Specify it in your .mmrc file or run the program with -%c (--%s) option.\n"
#define MSG_MISSING_PAR_NAME "%s: missing name. Specify it with -%c (--%s) option.\n"
#define MSG_MISSING_PAR_ID "%s: missing id. Specify it with -%c (--%s) option.\n"
#define MSG_ACCOUNT_NOT_FOUND "%s: given account not found or the name is ambiguous: %s.\n"
#define MSG_CATEGORY_NOT_FOUND "%s: given category not found or the name is ambiguous: %s.\n"
#define MSG_MAINCATEGORY_NOT_FOUND "%s: given main category not found or the name is ambiguous %s.\n"
#define MSG_CURRENCY_CANNOT_CHANGE "%s: changing currency for existing account is impossible. Don't use -%c (--%s) option when editing account details.\n"
#define MSG_MISSING_PAR_CURRENCY "%s: missing currency. Specify it with -%c (--%s) option or modify config file (~/.mmrc).\n"
#define MSG_MISSING_PAR_ACCOUNT_NAME "%s: missing account. Specify it with -%c (--%s) option.\n"
#define MSG_MISSING_PAR_CURRENCY_TO "%s: missing second currency. Specify it with -%c (--%s) option.\n"
#define MSG_MISSING_PAR_EXCHANGE_RATE "%s: missing exchange rate for the currencies. Specify it with -%c (--%s) option.\n"
#define MSG_MISSING_PAR_DESCRIPTION "%s: missing description. Specify it with -%c (--%s) option.\n"
#define MSG_MISSING_PAR_VALUE "%s: missing value. Specify it with -%c (--%s) option.\n"
#define MSG_MISSING_PAR_CATEGORY "%s: missing category. Specify it with -%c (--%s) option.\n"
#define MSG_MISSING_PAR_MAINCATEGORY "%s: missing main category. Specify it with -%c (--%s) option.\n"
#define MSG_MISSING_PAR_DATE "%s: missing date. Specify it with -%c (--%s) option.\n"
#define MSG_WRONG_DATE_FULL "%s: wrong date given. Specify it in format: YYYY-MM-DD.\n"
#define MSG_WRONG_DATE_MONTH "%s: wrong month given. Specify it in format: YYYY-MM.\n"
#define MSG_WRONG_DATE_MONTH_OR_YEAR "%s: wrong date given - expected year or month. Specify it in format: YYYY or YYYY-MM.\n"
#define MSG_MISSING_EXCHANGE_RATES "%s: currencies exchange rate(s) missing: %s - add them first.\n"
#define MSG_WRONG_OBJECT "%s: wrong object given.\n"

/* CONSTANTS */
#define FILE_PATH_MAX 1000
#define DATE_FULL_LEN 11

#define PAR_PROGNAME_LEN 11
#define PAR_NAME_LEN 11
#define PAR_DESCRIPTION_LEN 31
#define PAR_INSTITUTION_LEN 21
#define PAR_CURRENCY_LEN 4
#define PAR_VALUE_LEN 11
#define PAR_ID_NOT_SET -1

#define OBJ_OR_TYPE_LEN 20

#define NULL_STRING '\0'

/* FORMATTING STRINGS */
#define FS_GAP "  "
#define FS_GAPS " "
#define FS_ID "%5d"
#define FS_ID_T "%5s"
#define FS_NAME "%-10.10s"
#define FS_NAME_T "%-10s"
#define FS_ATYPE "%-10.10s"
#define FS_ATYPE_T "%-10s"
#define FS_CUR "%-3.3s"
#define FS_CUR_T "%-3s"
#define FS_CURL "%-13s"
#define FS_EXCHRT "%13.4f"
#define FS_EXCHRT_T "%-13s"
#define FS_INST "%-20.20s"
#define FS_INST_T "%-20s"
#define FS_DESC "%-30.30s"
#define FS_DESC_T "%-30s"
#define FS_CTYPE "%-11.11s"
#define FS_CTYPE_T "%-11s"
#define FS_MTYPE "%-8.8s"
#define FS_MTYPE_T "%-8s"
#define FS_DATE "%4d-%02d-%02d"
#define FS_DATE_T "%-10s"
#define FS_VALUE "%10.2f"
#define FS_VALUE_T "%10s"
#define FS_MONTH "%4d-%02d"
#define FS_MONTH_T "%-7s"


/* COMMANDS AND PARAMETERS */
#define OPTION_CMND_INIT_LONG "init"
#define OPTION_CMND_INIT_SHORT 'I'
#define OPTION_CMND_ADD_LONG "add"
#define OPTION_CMND_ADD_SHORT 'A'
#define OPTION_CMND_EDIT_LONG "edit"
#define OPTION_CMND_EDIT_SHORT 'E'
#define OPTION_CMND_DELETE_LONG "delete"
#define OPTION_CMND_DELETE_SHORT 'D'
#define OPTION_CMND_LIST_LONG "list"
#define OPTION_CMND_LIST_SHORT 'L'
#define OPTION_CMND_REPORT_LONG "report"
#define OPTION_CMND_REPORT_SHORT 'R'
#define OPTION_CMND_HELP_LONG "help"
#define OPTION_CMND_HELP_SHORT 'h'

#define OPTION_FILE_LONG "file"
#define OPTION_FILE_SHORT 'f'
#define OPTION_ID_LONG "id"
#define OPTION_ID_SHORT 'i'
#define OPTION_NAME_LONG "name"
#define OPTION_NAME_SHORT 'n'
#define OPTION_DESCRIPTION_LONG "description"
#define OPTION_DESCRIPTION_SHORT 's'
#define OPTION_BANK_LONG "bank"
#define OPTION_BANK_SHORT 'b'
#define OPTION_CURRENCY_LONG "currency"
#define OPTION_CURRENCY_SHORT 'j'
#define OPTION_CURRENCY_TO_LONG "currency-to"
#define OPTION_CURRENCY_TO_SHORT 'k'
#define OPTION_ACCOUNT_LONG "account"
#define OPTION_ACCOUNT_SHORT 'a'
#define OPTION_CATEGORY_LONG "category"
#define OPTION_CATEGORY_SHORT 'c'
#define OPTION_MAINCATEGORY_LONG "main-category"
#define OPTION_MAINCATEGORY_SHORT 'm'
#define OPTION_VALUE_LONG "value"
#define OPTION_VALUE_SHORT 'v'
#define OPTION_ACCOUNTTYPE_LONG "account-type"
#define OPTION_ACCOUNTTYPE_SHORT 'p'
#define OPTION_MAINCATEGORYTYPE_LONG "main-category-type"
#define OPTION_MAINCATEGORYTYPE_SHORT 'o'
#define OPTION_DATE_LONG "date"
#define OPTION_DATE_SHORT 'd'
#define OPTION_VERBOSE_LONG "verbose"
#define OPTION_VERBOSE_SHORT 1001

#define OPTIONS_LIST "IA:E:D:L:R:hf:i:n:s:b:j:k:a:c:m:v:o:p:d:"


/* DATA STRUCTURES */

/*
 * Enum with different types of accounts.
 */
typedef enum accountTypesT {
    ACC_TYPE_UNKNOWN = -1,
    ACC_TYPE_UNSET = 0,
    ACC_TRANSACTIONAL = 1,
    ACC_SAVING = 2,
    ACC_PROPERTY = 3,
    ACC_INVESTMENT = 4,
    ACC_LOAN = 5
} ACCOUNT_TYPE;

/*
 * Enum with different item (account, category etc.) statuses.
 */
typedef enum itemStatT {
    ITEM_STAT_CLOSED = 0,
    ITEM_STAT_OPEN = 1
} ITEM_STATUS;

/*
 * Enum with different category types.
 */
typedef enum maincategoryTypesT {
    CAT_TYPE_UNKNOWN = -3,
    CAT_TYPE_NOTSET = -2,
    CAT_COST = -1,
    CAT_TRANSFER = 0,
    CAT_INCOME = 1
} MAINCATEGORY_TYPE;

/**
 * Enum with different date type: to verify what part of date was entered by user
 */
typedef enum dateTypeT {
    DT_NO_DATE = 0,
    DT_YEAR = 1,
    DT_MONTH = 2,
    DT_FULL_DATE = 3
} DATE_TYPE;

/*
 * Struct with runtime parameters.
 * Contains application name, user parameters and configuration settings.
 */
typedef struct parametersT {
    char prog_name[PAR_PROGNAME_LEN];
    char dataFilePath[FILE_PATH_MAX];
    char name[PAR_NAME_LEN];
    char description[PAR_DESCRIPTION_LEN];
    char institution[PAR_INSTITUTION_LEN];
    ACCOUNT_TYPE account_type;
    char currency[PAR_CURRENCY_LEN];
    char currency_to[PAR_CURRENCY_LEN];
    char default_currency[PAR_CURRENCY_LEN];
    char acc_name[PAR_NAME_LEN];
    char cat_name[PAR_NAME_LEN];
    char maincat_name[PAR_NAME_LEN];
    MAINCATEGORY_TYPE maincategory_type;
    char value[PAR_VALUE_LEN];
    char date[DATE_FULL_LEN];
    char date_default[DATE_FULL_LEN];
    unsigned long id;           // if negative then it's not set
    int verbose;
} PARAMETERS;


/* FUNCTION PROTOTYPES */


/**
 * This function returns account id for given (part of) account name.
 * If there is no account or given name is ambiguous then the function returns -1.
 * @param acc_name char* with (part of) account name.
 * @return long id for found account.
 */
long account_id_for_name(sqlite3 * db, char *acc_name);

/**
 * Function returns ACCOUNT_TYPE number for given text.
 * @param account_type *char with account type
 * @return ACCOUNT_TYPE number.
 */
ACCOUNT_TYPE account_type_id(char *account_type);

/**
 * Function returns text description for account type.
 * @param acc_type ACC_TYPE enum value for which the text should be evaluated.
 * @return char* with text represenation of account type.
 */
char *account_type_text(ACCOUNT_TYPE acc_type);

/**
 * This function returns multiplicator for transaction value depending on category id.
 * @param db sqlite3* pointer to database
 * @param cat_id INT with category ID
 * @return INT: -1 for cost, 1 for income, 0 for transfer
 */
int category_factor_for_id(sqlite3 * db, int cat_id);


/**
 * This function returns category id for given (part of) category name.
 * If there is no category or a given name is ambiguous then the function returns -1.
 * @param cat_name char* with (part of) category name.
 * @return long id for found category.
 */
long category_id_for_name(sqlite3 * db, char *cat_name);

/**
 * Parses string with date (format: YYYY-MM-DD) and writes the values
 * to the following variables: year, month, day.
 * @param date_string *char with date
 * @param year_holder *int with year
 * @param month_holder *int with month
 * @param day_holder *int with day
 * @return DATA_TYPE with information about returned date
 * (0 - no date, 1 - only year, 2 - year-month, 3 - full date.
 */
DATE_TYPE date_from_string(char *date_string, int *year_holder,
                           int *month_holder, int *day_holder);

/**
 * Gets today date into date_holder.
 * @param date_holder char* to store the today date.
 */
void get_today(char *date_holder);

/**
 * Function returns text description for status of an item.
 * @param item_status ITEM_STATUS enum value for which the text should be evaluated.
 * @return char* with text representation of item status.
 */
char *item_status_text(ITEM_STATUS item_status);

/**
 * This function returns main category id for given (part of) its name.
 * If there is no main category or given name is ambiguous then the function returns -1.
 * @param main_category_name char* with (part of) main category name.
 * @return long id for found main category.
 */
long maincategory_id_for_name(sqlite3 * db, char *maincategory_name);

/**
 * Function returns multiplicator for given type.
 * To be used to correct values stored in accounts.
 * @param maincategory_type MAIN CATEGORY TYPE for which the factor is looked for
 * @return INT (-1) for cost, (1) for income, (1) for transfer
 */
int maincategory_type_factor(MAINCATEGORY_TYPE maincategory_type);

/**
 * Function returns CATEGORY_TYPE number for given text.
 * @param category_type *char with category type
 * @return CATEGORY_TYPE number.
 */
MAINCATEGORY_TYPE maincategory_type_id(char *maincategory_type);

/**
 * Function returns text description for category type
 * @param category_type CATEGORY_TYPE enum value for which the text should be evaluated.
 * @return char* with text representation of category type.
 */
char *maincategory_type_text(MAINCATEGORY_TYPE maincategory_type);

/**
 * Returns category id for given transaction
 * @param db sqlite3* database pointer
 * @param transaction_id unsigned long with id of the transaction
 * @return
 */
unsigned long transaction_category_id(sqlite3 * db,
                                      unsigned long transaction_id);


/* END OF INCLUDE GUARD */
#endif
