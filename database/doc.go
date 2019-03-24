/*
Package database handles the database operations.

The package exposes only a few of the underlying database/sql API, to force you
to use contexts and perform every operation inside a transaction.

The dbmodel_model.go defines all the type-safe CRUD operations. It is generated
by the generator at model/protoc-gen-dbmodel.

Custom queries should always be a SELECT statement and select only the primary
key column of some table. The application then fetchs the actual rows using the
CRUD operations. A custom query should never fetch values directly.
*/
package database
