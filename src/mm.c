/*
  Written 2015 by Marcin 'Zbroju' Zbroinski.
  Use of this source code is governed by a GNU General Public License
  that can be found in the LICENSE file.
*/

/* INCLUDES */
#include "mm.h"
#include "common.h"
#include "operations.h"
#include "reports.h"
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <getopt.h>
#include <libconfig.h>


/* CONSTANTS */

#define OBJECT_ACCOUNT_LONG "account"
#define OBJECT_ACCOUNT_SHORT "a"
#define OBJECT_TRANSACTION_LONG "transaction"
#define OBJECT_TRANSACTION_SHORT "t"
#define OBJECT_MAINCATEGORY_LONG "main-category"
#define OBJECT_MAINCATEGORY_SHORT "m"
#define OBJECT_CURRENCY_LONG "currency"
#define OBJECT_CURRENCY_SHORT "j"
#define OBJECT_CATEGORY_LONG "category"
#define OBJECT_CATEGORY_SHORT "c"
#define OBJECT_BUDGET_LONG "budget"
#define OBJECT_BUDGET_SHORT "b"
#define OBJECT_ACCOUNTSBALANCE_LONG "accounts-balance"
#define OBJECT_ACCOUNTSBALANCE_SHORT "ab"
#define OBJECT_ASSETSSUMMARY_LONG "assets-summary"
#define OBJECT_ASSETSSUMMARY_SHORT "as"
#define OBJECT_TRANSACTIONSBALANCE_LONG "transactions-balance"
#define OBJECT_TRANSACTIONSBALANCE_SHORT "tb"
#define OBJECT_CATEGORIESBALANCE_LONG "categories-balance"
#define OBJECT_CATEGORIESBALANCE_SHORT "cb"
#define OBJECT_MAINCATSBALANCE_LONG "main-categories-balance"
#define OBJECT_MAINCATSBALANCE_SHORT "mcb"
#define OBJECT_BUDGETCATS_LONG "budget-categories"
#define OBJECT_BUDGETCATS_SHORT "bc"
#define OBJECT_BUDGETMAINCATS_LONG "budget-main-categories"
#define OBJECT_BUDGETMAINCATS_SHORT "bmc"
#define OBJECT_NETVALUE_LONG "net-value"
#define OBJECT_NETVALUE_SHORT "nv"

/* DATA STRUCTURES */
typedef enum commandT {
    NO_COMMAND,
    INIT,
    ADD,
    EDIT,
    DELETE,
    LIST,
    REPORT,
    HELP
} COMMAND;

typedef enum objectT {
    NO_OBJECT,
    ACCOUNT,
    TRANSACTION,
    MAIN_CATEGORY,
    CATEGORY,
    CURRENCY,
    BUDGET,
    REP_ACCOUNTS_BALANCE,
    REP_ASSETS_SUMMARY,
    REP_TRANSACTIONS_BALANCE,
    REP_CATEGORIES_BALANCE,
    REP_MAINCATS_BALANCE,
    REP_BUDGET_CATS,
    REP_BUDGET_MAINCATS,
    REP_NETVALUE
} OBJECT;


/* FUNCTION PROTOTYPES */
/**
 * Prints error message when more than one command has been put by user.
 * @param app_name char* containig program name.
 * @return int with error code.
 */

/**
 * Assigns correct object for command.
 * @param arg char* object given by user.
 * @param obj OBJECT* variable which will hold the object.
 * @param app_name char* this program name.
 * @return int with error code.
 */
static int getObject(char *arg, OBJECT * obj, char *app_name);

/**
 * Reads config file and assigns config values to respective parameters.
 * @param parameters PARAMETERS to be set.
 * @return int with error code.
 */
static int getParametersFromConfFile(PARAMETERS * parameters);

/**
 * Prints errer message if there are too many commands given.
 */
static int perrorTooManyCmnds(char *app_name);

/**
 * Prints full help on standard output.
 */
static void print_help(void);

/**
 * Prints short summary of options on standard output.
 */
static void print_usage(void);

/**
 * Sets default parameters' values.
 * @param parameters PARAMETERS to be set.
 * @param app_name string with the application name.
 */
static void setDefaultParameters(PARAMETERS * parameters, char *app_name);


/* MAIN FUNCTION */
int main(int argc, char *const argv[])
{

    PARAMETERS parameters;
    COMMAND command = NO_COMMAND;
    OBJECT object = NO_OBJECT;
    int result = 0;

    setDefaultParameters(&parameters, argv[0]);

    if (getParametersFromConfFile(&parameters) != 0) {
        exit(EXIT_FAILURE);
    }
    // Define and parse user commands, objects and options
    static struct option cmdl_options[] = {
        {OPTION_CMND_INIT_LONG, no_argument, 0, OPTION_CMND_INIT_SHORT},
        {OPTION_CMND_ADD_LONG, required_argument, 0,
         OPTION_CMND_ADD_SHORT},
        {OPTION_CMND_EDIT_LONG, required_argument, 0,
         OPTION_CMND_EDIT_SHORT},
        {OPTION_CMND_DELETE_LONG, required_argument, 0,
         OPTION_CMND_DELETE_SHORT},
        {OPTION_CMND_LIST_LONG, required_argument, 0,
         OPTION_CMND_LIST_SHORT},
        {OPTION_CMND_REPORT_LONG, required_argument, 0,
         OPTION_CMND_REPORT_SHORT},
        {OPTION_CMND_HELP_LONG, no_argument, 0, OPTION_CMND_HELP_SHORT},
        {OPTION_FILE_LONG, required_argument, 0, OPTION_FILE_SHORT},
        {OPTION_ID_LONG, required_argument, 0, OPTION_ID_SHORT},
        {OPTION_NAME_LONG, required_argument, 0, OPTION_NAME_SHORT},
        {OPTION_DESCRIPTION_LONG, required_argument, 0,
         OPTION_DESCRIPTION_SHORT},
        {OPTION_BANK_LONG, required_argument, 0, OPTION_BANK_SHORT},
        {OPTION_CURRENCY_LONG, required_argument, 0,
         OPTION_CURRENCY_SHORT},
        {OPTION_CURRENCY_TO_LONG, required_argument, 0,
         OPTION_CURRENCY_TO_SHORT},
        {OPTION_ACCOUNT_LONG, required_argument, 0, OPTION_ACCOUNT_SHORT},
        {OPTION_CATEGORY_LONG, required_argument, 0,
         OPTION_CATEGORY_SHORT},
        {OPTION_MAINCATEGORY_LONG, required_argument, 0,
         OPTION_MAINCATEGORY_SHORT},
        {OPTION_VALUE_LONG, required_argument, 0, OPTION_VALUE_SHORT},
        {OPTION_ACCOUNTTYPE_LONG, required_argument, 0,
         OPTION_ACCOUNTTYPE_SHORT},
        {OPTION_MAINCATEGORYTYPE_LONG, required_argument, 0,
         OPTION_MAINCATEGORYTYPE_SHORT},
        {OPTION_DATE_LONG, required_argument, 0, OPTION_DATE_SHORT},
        {OPTION_VERBOSE_LONG, no_argument, 0, OPTION_VERBOSE_SHORT},
        {0, 0, 0, 0}
    };

    int opt = 0, long_index = 0;

    while ((opt =
            getopt_long(argc, argv, OPTIONS_LIST, cmdl_options,
                        &long_index)) != -1) {
        switch (opt) {
        case OPTION_CMND_INIT_SHORT:
            if (command == NO_COMMAND) {
                command = INIT;
            } else {
                result = perrorTooManyCmnds(parameters.prog_name);
            }
            break;
        case OPTION_CMND_ADD_SHORT:
            if (command == NO_COMMAND) {
                command = ADD;
                result = getObject(optarg, &object, parameters.prog_name);
            } else {
                result = perrorTooManyCmnds(parameters.prog_name);
            }
            break;
        case OPTION_CMND_EDIT_SHORT:
            if (command == NO_COMMAND) {
                command = EDIT;
                result = getObject(optarg, &object, parameters.prog_name);
            } else {
                result = perrorTooManyCmnds(parameters.prog_name);
            }
            break;
        case OPTION_CMND_DELETE_SHORT:
            if (command == NO_COMMAND) {
                command = DELETE;
                result = getObject(optarg, &object, parameters.prog_name);
            } else {
                result = perrorTooManyCmnds(parameters.prog_name);
            }
            break;
        case OPTION_CMND_LIST_SHORT:
            if (command == NO_COMMAND) {
                command = LIST;
                result = getObject(optarg, &object, parameters.prog_name);
            } else {
                result = perrorTooManyCmnds(parameters.prog_name);
            }
            break;
        case OPTION_CMND_REPORT_SHORT:
            if (command == NO_COMMAND) {
                command = REPORT;
                result = getObject(optarg, &object, parameters.prog_name);
            } else {
                result = perrorTooManyCmnds(parameters.prog_name);
            }
            break;
        case OPTION_CMND_HELP_SHORT:
            if (command == NO_COMMAND) {
                command = HELP;
            } else {
                result = perrorTooManyCmnds(parameters.prog_name);
            }
            break;
        case OPTION_FILE_SHORT:
            strncpy(parameters.dataFilePath, optarg, FILE_PATH_MAX - 1);
            break;
        case OPTION_ID_SHORT:
            parameters.id = atoi(optarg);
            break;
        case OPTION_NAME_SHORT:
            strncpy(parameters.name, optarg, PAR_NAME_LEN - 1);
            break;
        case OPTION_DESCRIPTION_SHORT:
            strncpy(parameters.description, optarg,
                    PAR_DESCRIPTION_LEN - 1);
            break;
        case OPTION_BANK_SHORT:
            strncpy(parameters.institution, optarg,
                    PAR_INSTITUTION_LEN - 1);
            break;
        case OPTION_CURRENCY_SHORT:
            strncpy(parameters.currency, optarg, PAR_CURRENCY_LEN - 1);
            break;
        case OPTION_CURRENCY_TO_SHORT:
            strncpy(parameters.currency_to, optarg, PAR_CURRENCY_LEN - 1);
            break;
        case OPTION_ACCOUNT_SHORT:
            strncpy(parameters.acc_name, optarg, PAR_NAME_LEN - 1);
            break;
        case OPTION_CATEGORY_SHORT:
            strncpy(parameters.cat_name, optarg, PAR_NAME_LEN - 1);
            break;
        case OPTION_MAINCATEGORY_SHORT:
            strncpy(parameters.maincat_name, optarg, PAR_NAME_LEN - 1);
            break;
        case OPTION_VALUE_SHORT:
            strncpy(parameters.value, optarg, PAR_VALUE_LEN - 1);
            break;
        case OPTION_ACCOUNTTYPE_SHORT:
            if ((parameters.account_type =
                 account_type_id(optarg)) == ACC_TYPE_UNKNOWN) {
                fprintf(stderr, "%s: uknown account type: %s.\n",
                        parameters.prog_name, optarg);
                result = 1;
            }
            break;
        case OPTION_MAINCATEGORYTYPE_SHORT:
            if ((parameters.maincategory_type =
                 maincategory_type_id(optarg)) == CAT_TYPE_UNKNOWN) {
                fprintf(stderr, "%s: unknown main category type: %s.\n",
                        parameters.prog_name, optarg);
                result = 1;
            }
            break;
        case OPTION_DATE_SHORT:
            strncpy(parameters.date, optarg, DATE_FULL_LEN - 1);
            break;
        case OPTION_VERBOSE_SHORT:
            parameters.verbose = 1;
            break;
        case '?':
            result = 1;
            break;
        default:
            fprintf(stderr, "%s: unknown option %c.\n",
                    parameters.prog_name, (char) opt);
            result = 1;
        }
    }

    if (result) {
        exit(EXIT_FAILURE);
    }
    // Display unknown options
    if (optind < argc) {
        fprintf(stderr, "%s: unknown options:", parameters.prog_name);
        while (optind < argc) {
            fprintf(stderr, " %s", argv[optind++]);
        }
        fprintf(stderr, ": Skipped.\n");
    }
    // Assign a function to respective command and its argument (object)
    if (command == NO_COMMAND) {
        print_usage();
        exit(EXIT_SUCCESS);
    } else if (command == INIT) {
        result = datafile_init(parameters);
    } else if (command == ADD) {
        switch (object) {
        case ACCOUNT:
            result = account_add(parameters);
            break;
        case TRANSACTION:
            result = transaction_add(parameters);
            break;
        case MAIN_CATEGORY:
            result = maincategory_add(parameters);
            break;
        case CATEGORY:
            result = category_add(parameters);
            break;
        case CURRENCY:
            result = currency_add(parameters);
            break;
        case BUDGET:
            result = budget_add(parameters);
            break;
        default:
            fprintf(stderr, MSG_WRONG_OBJECT, parameters.prog_name);
            result = EXIT_FAILURE;
        }
    } else if (command == EDIT) {
        switch (object) {
        case ACCOUNT:
            result = account_edit(parameters);
            break;
        case TRANSACTION:
            result = transaction_edit(parameters);
            break;
        case MAIN_CATEGORY:
            result = maincategory_edit(parameters);
            break;
        case CATEGORY:
            result = category_edit(parameters);
            break;
        case CURRENCY:
            result = currency_edit(parameters);
            break;
        case BUDGET:
            result = budget_edit(parameters);
            break;
        default:
            fprintf(stderr, MSG_WRONG_OBJECT, parameters.prog_name);
            result = EXIT_FAILURE;
        }
    } else if (command == DELETE) {
        switch (object) {
        case ACCOUNT:
            result = account_close(parameters);
            break;
        case TRANSACTION:
            result = transaction_remove(parameters);
            break;
        case MAIN_CATEGORY:
            result = maincategory_remove(parameters);
            break;
        case CATEGORY:
            result = category_remove(parameters);
            break;
        case CURRENCY:
            result = currency_remove(parameters);
            break;
        case BUDGET:
            result = budget_remove(parameters);
            break;
        default:
            fprintf(stderr, MSG_WRONG_OBJECT, parameters.prog_name);
            result = EXIT_FAILURE;
        }
    } else if (command == LIST) {
        switch (object) {
        case ACCOUNT:
            result = account_list(parameters);
            break;
        case TRANSACTION:
            result = transaction_list(parameters);
            break;
        case MAIN_CATEGORY:
            result = maincategory_list(parameters);
            break;
        case CATEGORY:
            result = category_list(parameters);
            break;
        case CURRENCY:
            result = currency_list(parameters);
            break;
        case BUDGET:
            result = budget_list(parameters);
            break;
        default:
            fprintf(stderr, MSG_WRONG_OBJECT, parameters.prog_name);
            result = EXIT_FAILURE;
        }
    } else if (command == REPORT) {
        switch (object) {
        case REP_ACCOUNTS_BALANCE:
            result = accounts_balance(parameters);
            break;
        case REP_ASSETS_SUMMARY:
            result = assets_summary(parameters);
            break;
        case REP_TRANSACTIONS_BALANCE:
            result = transactions_balance(parameters);
            break;
        case REP_CATEGORIES_BALANCE:
            result = categories_balance(parameters);
            break;
        case REP_MAINCATS_BALANCE:
            result = maincategories_balance(parameters);
            break;
        case REP_BUDGET_CATS:
            result = budget_report_categories(parameters);
            break;
        case REP_BUDGET_MAINCATS:
            result = budget_report_maincategories(parameters);
            break;
        case REP_NETVALUE:
            result = net_value(parameters);
            break;
        default:
            fprintf(stderr, MSG_WRONG_OBJECT, parameters.prog_name);
            result = EXIT_FAILURE;
        }
    } else if (command == HELP) {
        print_help();
        result = EXIT_SUCCESS;
    }

    exit(result);
}

/* SUPPORTIVE FUNCTIONS */

static int getObject(char *arg, OBJECT * obj, char *app_name)
{
    int result = 0;

    if (*obj == NO_OBJECT) {
        if (strncmp(arg, OBJECT_ACCOUNT_LONG, OBJ_OR_TYPE_LEN) == 0
            || strncmp(arg, OBJECT_ACCOUNT_SHORT, OBJ_OR_TYPE_LEN) == 0) {
            *obj = ACCOUNT;
        } else if (strncmp(arg, OBJECT_TRANSACTION_LONG, OBJ_OR_TYPE_LEN)
                   == 0
                   || strncmp(arg, OBJECT_TRANSACTION_SHORT,
                              OBJ_OR_TYPE_LEN) == 0) {
            *obj = TRANSACTION;
        } else if (strncmp(arg, OBJECT_MAINCATEGORY_LONG, OBJ_OR_TYPE_LEN)
                   == 0
                   || strncmp(arg, OBJECT_MAINCATEGORY_SHORT,
                              OBJ_OR_TYPE_LEN) == 0) {
            *obj = MAIN_CATEGORY;
        } else if (strncmp(arg, OBJECT_CATEGORY_LONG, OBJ_OR_TYPE_LEN) == 0
                   || strncmp(arg, OBJECT_CATEGORY_SHORT,
                              OBJ_OR_TYPE_LEN) == 0) {
            *obj = CATEGORY;
        } else if (strncmp(arg, OBJECT_CURRENCY_LONG, OBJ_OR_TYPE_LEN) == 0
                   || strncmp(arg, OBJECT_CURRENCY_SHORT,
                              OBJ_OR_TYPE_LEN) == 0) {
            *obj = CURRENCY;
        } else if (strncmp(arg, OBJECT_BUDGET_LONG, OBJ_OR_TYPE_LEN) == 0
                   || strncmp(arg, OBJECT_BUDGET_SHORT,
                              OBJ_OR_TYPE_LEN) == 0) {
            *obj = BUDGET;
        } else
            if (strncmp(arg, OBJECT_ACCOUNTSBALANCE_LONG, OBJ_OR_TYPE_LEN)
                == 0
                || strncmp(arg, OBJECT_ACCOUNTSBALANCE_SHORT,
                           OBJ_OR_TYPE_LEN) == 0) {
            *obj = REP_ACCOUNTS_BALANCE;
        } else
            if (strncmp
                (arg, OBJECT_TRANSACTIONSBALANCE_LONG, OBJ_OR_TYPE_LEN)
                == 0
                || strncmp(arg, OBJECT_TRANSACTIONSBALANCE_SHORT,
                           OBJ_OR_TYPE_LEN) == 0) {
            *obj = REP_TRANSACTIONS_BALANCE;
        } else if (strncmp(arg, OBJECT_ASSETSSUMMARY_LONG, OBJ_OR_TYPE_LEN)
                   == 0
                   || strncmp(arg, OBJECT_ASSETSSUMMARY_SHORT,
                              OBJ_OR_TYPE_LEN) == 0) {
            *obj = REP_ASSETS_SUMMARY;
        } else
            if (strncmp
                (arg, OBJECT_CATEGORIESBALANCE_LONG, OBJ_OR_TYPE_LEN)
                == 0
                || strncmp(arg, OBJECT_CATEGORIESBALANCE_SHORT,
                           OBJ_OR_TYPE_LEN) == 0) {
            *obj = REP_CATEGORIES_BALANCE;
        } else
            if (strncmp(arg, OBJECT_MAINCATSBALANCE_LONG, OBJ_OR_TYPE_LEN)
                == 0
                || strncmp(arg, OBJECT_MAINCATSBALANCE_SHORT,
                           OBJ_OR_TYPE_LEN) == 0) {
            *obj = REP_MAINCATS_BALANCE;
        } else if (strncmp(arg, OBJECT_BUDGETCATS_LONG, OBJ_OR_TYPE_LEN) ==
                   0
                   || strncmp(arg, OBJECT_BUDGETCATS_SHORT,
                              OBJ_OR_TYPE_LEN) == 0) {
            *obj = REP_BUDGET_CATS;
        } else
            if (strncmp(arg, OBJECT_BUDGETMAINCATS_LONG, OBJ_OR_TYPE_LEN)
                == 0
                || strncmp(arg, OBJECT_BUDGETMAINCATS_SHORT,
                           OBJ_OR_TYPE_LEN) == 0) {
            *obj = REP_BUDGET_MAINCATS;
        } else if (strncmp(arg, OBJECT_NETVALUE_LONG, OBJ_OR_TYPE_LEN) == 0
                   || strncmp(arg, OBJECT_NETVALUE_SHORT,
                              OBJ_OR_TYPE_LEN) == 0) {
            *obj = REP_NETVALUE;
        } else {
            fprintf(stderr, "%s: unknown <object> - %s.\n", app_name, arg);
            result = 1;
        }
    } else {
        fprintf(stderr, "%s: more than one object given - %s. Stop.\n",
                app_name, arg);
        result = 1;
    }
    return result;
}

static int getParametersFromConfFile(PARAMETERS * parameters)
{
    char conf_file_path[FILE_PATH_MAX] = { NULL_STRING };
    config_t cfg;
    const char *str;

    strcpy(conf_file_path, getenv("HOME"));
    strcat(conf_file_path, "/.mmrc");

    // Read the file. If failure - report & exit.
    config_init(&cfg);
    if (!config_read_file(&cfg, conf_file_path)) {
        fprintf(stderr, "%s: %s in file %s (line: %d)\n",
                parameters->prog_name, config_error_text(&cfg),
                config_error_file(&cfg), config_error_line(&cfg));
        config_destroy(&cfg);
        return (EXIT_FAILURE);
    }
    // Get the DATA_FILE
    if (config_lookup_string(&cfg, "DATA_FILE", &str)) {
        strncpy(parameters->dataFilePath, str, FILE_PATH_MAX - 1);
    }
    // Get the DEFAULT_CURRENCY
    if (config_lookup_string(&cfg, "DEFAULT_CURRENCY", &str)) {
        strncpy(parameters->default_currency, str, PAR_CURRENCY_LEN - 1);
    }
    // Get the VERBOSE flague
    int verbose_flague;
    if (config_lookup_bool(&cfg, "VERBOSE", &verbose_flague)) {
        parameters->verbose = verbose_flague;
    }

    config_destroy(&cfg);
    return (EXIT_SUCCESS);
}

static int perrorTooManyCmnds(char *app_name)
{
    fprintf(stderr, "%s: more than one command given. Stop.\n", app_name);
    return 1;
}

static void print_help(void)
{
    printf("Usage:\n");
    printf("\tmm COMMAND [object | reports] [OPTIONS]\n");
    printf("\nCOMMANDS:\n");
    printf("\t-%c, --%s\tinit a new file. Requires -%c (--%s) option.\n",
           OPTION_CMND_INIT_SHORT, OPTION_CMND_INIT_LONG,
           OPTION_FILE_SHORT, OPTION_FILE_LONG);
    printf("\t-%c, --%s\tadd new <object> to file.\n",
           OPTION_CMND_ADD_SHORT, OPTION_CMND_ADD_LONG);
    printf
        ("\t-%c, --%s\tedit existing <object>. Requires -%c (--%s) option to indicate the object.\n",
         OPTION_CMND_EDIT_SHORT, OPTION_CMND_EDIT_LONG, OPTION_ID_SHORT,
         OPTION_ID_LONG);
    printf
        ("\t-%c, --%s\tdelete existing <object>. Requires -%c (--%s) option to indicate the object.\n",
         OPTION_CMND_DELETE_SHORT, OPTION_CMND_DELETE_LONG,
         OPTION_ID_SHORT, OPTION_ID_LONG);
    printf
        ("\t-%c, --%s\tlist <objects>. You can apply filters for the <objects>.\n",
         OPTION_CMND_LIST_SHORT, OPTION_CMND_LIST_LONG);
    printf
        ("\t-%c, --%s\tshow <report>. You can apply filters for the <report>.\n",
         OPTION_CMND_REPORT_SHORT, OPTION_CMND_REPORT_LONG);
    printf("\t-%c, --%s\tshow this help information.\n",
           OPTION_CMND_HELP_SHORT, OPTION_CMND_HELP_LONG);
    printf("\nOBJECTS:\n");
    printf("\t%s, %s\tobject to manipulate accounts.\n",
           OBJECT_ACCOUNT_SHORT, OBJECT_ACCOUNT_LONG);
    printf("\t%s, %s\tobject to manipulate transactions.\n",
           OBJECT_TRANSACTION_SHORT, OBJECT_TRANSACTION_LONG);
    printf("\t%s, %s\tobject to manipulate main categories.\n",
           OBJECT_MAINCATEGORY_SHORT, OBJECT_MAINCATEGORY_LONG);
    printf("\t%s, %s\tobject to manipulate currencies.\n",
           OBJECT_CURRENCY_SHORT, OBJECT_CURRENCY_LONG);
    printf("\t%s, %s\tobject to manipulate categories.\n",
           OBJECT_CATEGORY_SHORT, OBJECT_CATEGORY_LONG);
    printf("\t%s, %s\tobject to manipulate budgets.\n",
           OBJECT_BUDGET_SHORT, OBJECT_BUDGET_LONG);
    printf("\nREPORTS:\n");
    printf("\t%s, %s\tobject to show report of accounts balances.\n",
           OBJECT_ACCOUNTSBALANCE_SHORT, OBJECT_ACCOUNTSBALANCE_LONG);
    printf("\t%s, %s\tobject to show report of assets summary.\n",
           OBJECT_ASSETSSUMMARY_SHORT, OBJECT_ASSETSSUMMARY_LONG);
    printf("\t%s, %s\tobject to show report of transactions balances.\n",
           OBJECT_TRANSACTIONSBALANCE_SHORT,
           OBJECT_TRANSACTIONSBALANCE_LONG);
    printf("\t%s, %s\tobject to show report of categories balances.\n",
           OBJECT_CATEGORIESBALANCE_SHORT, OBJECT_CATEGORIESBALANCE_LONG);
    printf
        ("\t%s, %s\tobject to show report of main categories balances.\n",
         OBJECT_MAINCATSBALANCE_SHORT, OBJECT_MAINCATSBALANCE_LONG);
    printf("\t%s, %s\tobject to show report of budget for categories.\n",
           OBJECT_BUDGETCATS_SHORT, OBJECT_BUDGETCATS_LONG);
    printf
        ("\t%s, %s\tobject to show report of budget for main categories.\n",
         OBJECT_BUDGETMAINCATS_SHORT, OBJECT_BUDGETMAINCATS_LONG);
    printf("\t%s, %s\tobject to show report of net value.\n",
           OBJECT_NETVALUE_SHORT, OBJECT_NETVALUE_LONG);
    printf("\nOPTIONS:\n");
    printf("\t-%c, --%s\tfull path to data file.\n", OPTION_FILE_SHORT,
           OPTION_FILE_LONG);
    printf("\t-%c, --%s\tid for identifying particular object.\n",
           OPTION_ID_SHORT, OPTION_ID_LONG);
    printf
        ("\t-%c, --%s\tname of an object (account, main category & category).\n",
         OPTION_NAME_SHORT, OPTION_NAME_LONG);
    printf("\t-%c, --%s\tdescription of a transaction.\n",
           OPTION_DESCRIPTION_SHORT, OPTION_DESCRIPTION_LONG);
    printf("\t-%c, --%s\tbank name where a given account is maintained.\n",
           OPTION_BANK_SHORT, OPTION_BANK_LONG);
    printf("\t-%c, --%s\tcurrency.\n", OPTION_CURRENCY_SHORT,
           OPTION_CURRENCY_LONG);
    printf("\t-%c, --%s\tcurrency against.\n", OPTION_CURRENCY_TO_SHORT,
           OPTION_CURRENCY_TO_LONG);
    printf
        ("\t-%c, --%s\taccount name. It's enough to give part of the name as long as it allows to identify one account.\n",
         OPTION_ACCOUNT_SHORT, OPTION_ACCOUNT_LONG);
    printf
        ("\t-%c, --%s\tcategory name. It's enough to give part of the name as long as it allows to identify one category.\n",
         OPTION_CATEGORY_SHORT, OPTION_CATEGORY_LONG);
    printf
        ("\t-%c, --%s\tmain category name. It's enough to give part of the name as long as it allows to identify one main category.\n",
         OPTION_MAINCATEGORY_SHORT, OPTION_MAINCATEGORY_LONG);
    printf
        ("\t-%c, --%s\tvalue of a transaction, or exchange rate when working with currency.\n",
         OPTION_VALUE_SHORT, OPTION_VALUE_LONG);
    printf
        ("\t-%c, --%s\taccount type. Allowed values are: t/transact (default), s/saving, p/property, i/investment, l/loan.\n",
         OPTION_ACCOUNTTYPE_SHORT, OPTION_ACCOUNTTYPE_LONG);
    printf
        ("\t-%c, --%s\tmain category type. Allowed values are: c/cost, t/transfer, i/income.\n",
         OPTION_MAINCATEGORYTYPE_SHORT, OPTION_MAINCATEGORYTYPE_LONG);
    printf
        ("\t-%c, --%s\tdate. Required format is YYYY (for year), YYYY-MM (for year-month) and YYYY-MM-DD (for full date). Today by default.\n",
         OPTION_DATE_SHORT, OPTION_DATE_LONG);
    printf("\t  , --%s\tmake the program verbose.\n", OPTION_VERBOSE_LONG);
}

static void print_usage(void)
{
    printf("Usage:\n");
    printf("\tmm COMMAND [object | reports] [OPTIONS] [--%s]\n",
           OPTION_VERBOSE_LONG);

    printf("\tmm -%c%c%c%c%c%c%c"
           " [%s %s %s %s %s %s | %s %s %s %s %s %s %s %s]"
           " [-%c%c%c%c%c%c%c%c%c%c%c%c%c%c]"
           " [--%s]\n", OPTION_CMND_INIT_SHORT, OPTION_CMND_ADD_SHORT,
           OPTION_CMND_EDIT_SHORT, OPTION_CMND_DELETE_SHORT,
           OPTION_CMND_LIST_SHORT, OPTION_CMND_REPORT_SHORT,
           OPTION_CMND_HELP_SHORT, OBJECT_ACCOUNT_SHORT,
           OBJECT_TRANSACTION_SHORT, OBJECT_MAINCATEGORY_SHORT,
           OBJECT_CURRENCY_SHORT, OBJECT_CATEGORY_SHORT,
           OBJECT_BUDGET_SHORT, OBJECT_ACCOUNTSBALANCE_SHORT,
           OBJECT_ASSETSSUMMARY_SHORT, OBJECT_TRANSACTIONSBALANCE_SHORT,
           OBJECT_CATEGORIESBALANCE_SHORT, OBJECT_MAINCATSBALANCE_SHORT,
           OBJECT_BUDGETCATS_SHORT, OBJECT_BUDGETMAINCATS_SHORT,
           OBJECT_NETVALUE_SHORT, OPTION_FILE_SHORT, OPTION_ID_SHORT,
           OPTION_NAME_SHORT, OPTION_DESCRIPTION_SHORT, OPTION_BANK_SHORT,
           OPTION_CURRENCY_SHORT, OPTION_CURRENCY_TO_SHORT,
           OPTION_ACCOUNT_SHORT, OPTION_CATEGORY_SHORT,
           OPTION_MAINCATEGORY_SHORT, OPTION_VALUE_SHORT,
           OPTION_ACCOUNTTYPE_SHORT, OPTION_MAINCATEGORYTYPE_SHORT,
           OPTION_DATE_SHORT, OPTION_VERBOSE_LONG);
}

static void setDefaultParameters(PARAMETERS * parameters, char *app_name)
{
    memset(parameters->prog_name, NULL_STRING, PAR_PROGNAME_LEN);
    memset(parameters->dataFilePath, NULL_STRING, FILE_PATH_MAX);
    memset(parameters->name, NULL_STRING, PAR_NAME_LEN);
    memset(parameters->description, NULL_STRING, PAR_DESCRIPTION_LEN);
    memset(parameters->institution, NULL_STRING, PAR_INSTITUTION_LEN);
    memset(parameters->currency, NULL_STRING, PAR_CURRENCY_LEN);
    memset(parameters->currency_to, NULL_STRING, PAR_CURRENCY_LEN);
    memset(parameters->default_currency, NULL_STRING, PAR_CURRENCY_LEN);
    memset(parameters->acc_name, NULL_STRING, PAR_NAME_LEN);
    memset(parameters->cat_name, NULL_STRING, PAR_NAME_LEN);
    memset(parameters->maincat_name, NULL_STRING, PAR_NAME_LEN);
    memset(parameters->value, NULL_STRING, PAR_VALUE_LEN);
    memset(parameters->date, NULL_STRING, DATE_FULL_LEN);
    memset(parameters->date_default, NULL_STRING, DATE_FULL_LEN);

    parameters->account_type = ACC_TYPE_UNSET;
    parameters->maincategory_type = CAT_TYPE_NOTSET;
    get_today(parameters->date_default);
    parameters->id = PAR_ID_NOT_SET;
    parameters->verbose = 0;
    strncpy(parameters->prog_name, app_name, PAR_PROGNAME_LEN - 1);
}
