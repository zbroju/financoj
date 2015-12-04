/*
 * File:   operations.h
 * Author: Marcin 'Zbroju' Zbroinski <marcin at zbroinski.net>
 *
 * Created on 10 kwiecień 2015, 19:27
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but without ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation , Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA
 */

#ifndef OPERATIONS_H
#define	OPERATIONS_H

/* INCLUDES */
#include "common.h"
#include <sqlite3.h>

/* FUNCTION PROTOTYPES */

/**
 * Adds a new account.
 * @param parameters PARAMETERS (struct) containing the new account data.
 * @return int with error code.
 */
int account_add(PARAMETERS parameters);

/**
 * Edit data of existing account.
 * @param parameters PARAMETERS (struct) containing the edited account data.
 * @return int with error code.
 */
int account_edit(PARAMETERS parameters);

/**
 * Remove account.
 * This funtctions actually does not remove the account, because transactions
 * assigned to it would be orphaned.
 * Instead, a special status is set for given account (CLOSED).
 * @param parameters PARAMETERS (struct) containing account id.
 * @return int with error code.
 */
int account_close(PARAMETERS parameters);

/**
 * List accounts.
 * @param parameters PARAMETERS (struct) containing commandline arguments.
 * @return int with error code.
 */
int account_list(PARAMETERS parameters);

/**
 * Adds a new category.
 * @param parameters PARAMETERS (struct) containing the new category data.
 * @return int with error code.
 */
int category_add(PARAMETERS parameters);

/**
 * Edit data of existing category.
 * @param parameters PARAMETERS (struct) containing the edited category data.
 * @return int with error code.
 */
int category_edit(PARAMETERS parameters);

/**
 * Remove category.
 * This funtctions actually does not remove the category, because
 * transactions assigned to it would be orphaned.
 * Instead, a special status is set for given category (CLOSED).
 * @param parameters PARAMETERS (struct) containing category id.
 * @return int with error code.
 */
int category_remove(PARAMETERS parameters);

/**
 * List categories.
 * @param parameters PARAMETERS (struct) containing commandline arguments.
 * @return int with error code.
 */
int category_list(PARAMETERS parameters);

/**
 * Adds a new currency exchange rate.
 * @param parameters PARAMETERS (struct) containing the new currency exchange rate data.
 * @return int with error code.
 */
int currency_add(PARAMETERS parameters);

/**
 * Edit data of existing currency exchange rate.
 * @param parameters PARAMETERS (struct) containing the edited currency exchange rate data.
 * @return int with error code.
 */
int currency_edit(PARAMETERS parameters);

/**
 * Remove exchange rate.
 * @param parameters PARAMETERS (struct) containing currencies data.
 * @return int with error code.
 */
int currency_remove(PARAMETERS parameters);

/**
 * List currencies exchange rates.
 * @param parameters PARAMETERS (struct) containing commandline arguments.
 * @return int with error code.
 */
int currency_list(PARAMETERS parameters);

/**
 * Creates database file.
 * The file is created with path and name given in PARAMETERS.
 * @param parameters PARAMETERS containg data file path.
 * @return int with error code.
 */
int datafile_init(PARAMETERS parameters);

/**
 * Adds a new main category.
 * @param parameters PARAMETERS (struct) containing the new main category data.
 * @return int with error code.
 */
int maincategory_add(PARAMETERS parameters);

/**
 * Edit data of existing main category.
 * @param parameters PARAMETERS (struct) containing the edited main category data.
 * @return int with error code.
 */
int maincategory_edit(PARAMETERS parameters);

/**
 * Remove main category.
 * This funtctions actually does not remove the main category, because
 * categories and thus transactions assigned to it would be orphaned.
 * Instead, a special status is set for given main category (CLOSED).
 * @param parameters PARAMETERS (struct) containing main category id.
 * @return int with error code.
 */
int maincategory_remove(PARAMETERS parameters);

/**
 * List main categories.
 * @param parameters PARAMETERS (struct) containing commandline arguments.
 * @return int with error code.
 */
int maincategory_list(PARAMETERS parameters);



/**
 * Adds a new transaction.
 * @param parameters PARAMETERS (struct) containing the new account data.
 * @return int with error code.
 */
int transaction_add(PARAMETERS parameters);

/**
 * Edit data of existing transaction.
 * @param parameters PARAMETERS (struct) containing the new account data.
 * @return int with error code.
 */
int transaction_edit(PARAMETERS parameters);

/**
 * Remove transaction.
 * @param parameters PARAMETERS (struct) containing transaction id.
 * @return int with error code.
 */
int transaction_remove(PARAMETERS parameters);

/**
 * List transactions.
 * @param parameters PARAMETERS (struct) containing commandline arguments.
 * @return int with error code.
 */
int transaction_list(PARAMETERS parameters);



/**
 *
 * @param parameters
 * @return
 */
int budget_add(PARAMETERS parameters);

/**
 * List budgets.
 * @param parameters PARAMETERS (struct) containing commandline arguments.
 * @return int with error code.
 */
int budget_list(PARAMETERS parameters);

/**
 * Edit data of existing budget
 * @param parameters PARAMETERS (struct) containing transaction id.
 * @return int with error code.
 */
int budget_edit(PARAMETERS parameters);

/**
 * Remove given budget
 * @param parameters PARAMETERS (struct) containing commandline arguments.
 * @return int with error code.
 */
int budget_remove(PARAMETERS parameters);

#endif	/* OPERATIONS_H */

