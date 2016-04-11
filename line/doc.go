/* Package line is used to define lexicon line format for parsing input and printing output.

Interfaces:
* Format - simple line format definition (field names and indices)
* Parser - a more complex parser, containing a Format definition, but also adds the possibility to write specific code for parsing that cannot be handeled by the Format specs alone (multi-value fields, etc).
*/
package line
