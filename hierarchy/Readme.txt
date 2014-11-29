PACKAGE DOCUMENTATION

package hierarchy
    import "github.com/dschnare/doctag/hierarchy"

    Package hierarchy implements a doctag transformer that transforms a list
    of doctags into a map hierarchy. The name of each doctag must use a
    separator character to describe a hierarchy. Doctags can also support
    adding to a slice by prefixing the doctag name with "#". If no separator
    character is found then the hierarchy will be flat.

    Example:

    Doctag document:

	<{ page/title }>This is the page title<{!}>
	<{ page/keywords }>awesome,stuff,aboutpeople<{!}>

    Map hierarchy:

	map{
	   "page": map{
	      "title": "This is the page title",
	      "keywords": "awesome,stuff,aboutpeople",
	   },
	}

    By prefixing doctags with the '#' character it indicates to the
    transformer that the implicitly created map (or string value) will be
    appended to a slice. Aany doctags that add keys to the implicitly
    created map that has been added to a slice will be set on the last map
    added to the slice.

    Example:

    Doctag document:

	<{ page/title }>This is the page title<{!}>
	<{ ! These doctags will append each string value to a slice indexed by "keywords" }>
	<{ page/#keywords }>awesome<{!}>
	<{ page/#keywords }>stuff<{!}>
	<{ page/#keywords }>aboutpeople<{!}>

	<{ page/content }>
	Some stuff about people

	<{ ! These doctags will append the implicitly created map to a slice indexed by "links" }>
	<{ page/#links/rel }>alternate<{!}>
	<{ page/links/href }>http://my.domain.com/alternate.html<{!}>
	<{ page/#links/rel }>next<{!}>
	<{ page/links/href }>http://my.domain.com/next.html<{!}>
	<{ page/#links/rel }>prev<{!}>
	<{ page/links/href }>http://my.domain.com/prev.html<{!}>

    Map hierarchy:

	map{
	   "page": map{
	      "title": "This is the page title",
	      "keywords": ["awesome", "stuff", "aboutpeople"],
	      "links": [
	        map{
	          "rel": "alternate",
	          "href": "http://my.domain.com/alternate.html",
	        },
	        map{
	          "rel": "next",
	          "href": "http://my.domain.com/next.html",
	        },
	        map{
	          "rel": "prev",
	          "href": "http://my.domain.com/prev.html",
	        },
	      ],
	   },
	}

CONSTANTS

const DefaultSeparator = '/'
    DefaultSeparator is a constant for the default character used to delimit
    separate doctag names.

FUNCTIONS

func Transform(doctags []*parse.DoctagNode, jsonKeysToIdentifiers bool) (map[string]interface{}, error)
    Transform transforms a slice of DoctagNodes into a hierarchical map that
    represents a JSON object. The default separater character will be used
    when parsing hierarchical doctags.

func TransformWithSeparator(doctags []*parse.DoctagNode, jsonKeysToIdentifiers bool, separator rune) (map[string]interface{}, error)
    TransformWithSeparator transforms a slice of DoctagNodes with a specific
    doctag separator character into a hierarchical map that represents a
    JSON object.

SUBDIRECTORIES

	fixtures

