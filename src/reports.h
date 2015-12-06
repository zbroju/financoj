/*
  Written 2015 by Marcin 'Zbroju' Zbroinski.
  Use of this source code is governed by a GNU General Public License
  that can be found in the LICENSE file.
*/

#ifndef REPORTS_H
#define	REPORTS_H

#include "common.h"

/**
 * Shows summary of all accounts in theirs original currencies.
 * The balance is calculated for a given date or today if the date is not specified.
 * @param parameters runtime parameters.
 * @return INT with error code.
 */
int accounts_balance(PARAMETERS parameters);

/**
 * Shows summary of all accounts in one currency and current net value.
 * @param parameters runtime parameters.
 * @return INT with error code.
 */
int assets_summary(PARAMETERS parameters);

/**
 * Shows transactions for given filters (date, main category, category, account)
 * @param parameters runtime parameters
 * @return INT with error code.
 */
int transactions_balance(PARAMETERS parameters);

/**
 * Shows summary of categories for given time and currency.
 * @param parameters runtime parameters
 * @return INT with error code.
 */
int categories_balance(PARAMETERS parameters);

/**
 * Shows summary of main categories for given time and currency.
 * @param parameters runtime parameters
 * @return INT with error code.
 */
int maincategories_balance(PARAMETERS parameters);

/**
 * Shows budget limits vs actual spent money on category level.
 * @param parameters runtime parameters
 * @return INT with error code.
 */
int budget_report_categories(PARAMETERS parameters);

/**
 * Shows budget limits vs actual spent money on main category level.
 * @param parameters runtime parameters
 * @return INT with error code.
 */
int budget_report_maincategories(PARAMETERS parameters);

/**
 * Shows historical net value.
 * @param parameters runtime parameters
 * @return INT with error code.
 */
int net_value(PARAMETERS parameters);


#endif                          /* REPORTS_H */
