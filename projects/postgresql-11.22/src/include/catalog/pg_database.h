/*-------------------------------------------------------------------------
 *
 * pg_database.h
 *	  definition of the "database" system catalog (pg_database)
 *
 *
 * Portions Copyright (c) 1996-2018, PostgreSQL Global Development Group
 * Portions Copyright (c) 1994, Regents of the University of California
 *
 * src/include/catalog/pg_database.h
 *
 * NOTES
 *	  The Catalog.pm module reads this file and derives schema
 *	  information.
 *
 *-------------------------------------------------------------------------
 */
#ifndef PG_DATABASE_H
#define PG_DATABASE_H

#include "catalog/genbki.h"
#include "catalog/pg_database_d.h"

/* ----------------
 *		pg_database definition.  cpp turns this into
 *		typedef struct FormData_pg_database
 * ----------------
 */
CATALOG(pg_database,1262,DatabaseRelationId) BKI_SHARED_RELATION BKI_ROWTYPE_OID(1248,DatabaseRelation_Rowtype_Id) BKI_SCHEMA_MACRO
{
	NameData	datname;		/* database name */
	Oid			datdba;			/* owner of database */
	int32		encoding;		/* character encoding */
	NameData	datcollate;		/* LC_COLLATE setting */
	NameData	datctype;		/* LC_CTYPE setting */
	bool		datistemplate;	/* allowed as CREATE DATABASE template? */
	bool		datallowconn;	/* new connections allowed? */
	/*
	 * Max connections allowed. Negative values have special meaning, see
	 * DATCONNLIMIT_* defines below.
	 */
	int32		datconnlimit;
	Oid			datlastsysoid;	/* highest OID to consider a system OID */
	TransactionId datfrozenxid; /* all Xids < this are frozen in this DB */
	TransactionId datminmxid;	/* all multixacts in the DB are >= this */
	Oid			dattablespace;	/* default table space for this DB */

#ifdef CATALOG_VARLEN			/* variable-length fields start here */
	aclitem		datacl[1];		/* access permissions */
#endif
} FormData_pg_database;

/* ----------------
 *		Form_pg_database corresponds to a pointer to a tuple with
 *		the format of pg_database relation.
 * ----------------
 */
typedef FormData_pg_database *Form_pg_database;

/*
 * Special values for pg_database.datconnlimit. Normal values are >= 0.
 */
#define		  DATCONNLIMIT_UNLIMITED	-1	/* no limit */

/*
 * A database is set to invalid partway through being dropped.  Using
 * datconnlimit=-2 for this purpose isn't particularly clean, but is
 * backpatchable.
 */
#define		  DATCONNLIMIT_INVALID_DB	-2

extern bool database_is_invalid_form(Form_pg_database datform);
extern bool database_is_invalid_oid(Oid dboid);

#endif							/* PG_DATABASE_H */
