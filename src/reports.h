/*
 * File:   reports.h
 * Author: Marcin 'Zbroju' Zbroinski <marcin at zbroinski.net>
 *
 * Created on 12 kwiecień 2015, 17:49
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


#endif	/* REPORTS_H */

