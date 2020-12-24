// gluamapper provides an easy way to map GopherLua values to Go values.
//
// Field Tags
//
// When mapping to a struct, gluamapper will use the field name by
// default to perform the mapping. For example, if a struct has a field
// "Username" then gluamapper will look for a key in the Lua table of "Username".
// You can use struct tags to look for a different key name in the Lua table.
// See example Mapper (tagName)
//
// Unexported fields
//
// Since unexported (private) struct fields cannot be set outside the package
// where they are defined, gluamapper will simply skip them.
//
package gluamapper
