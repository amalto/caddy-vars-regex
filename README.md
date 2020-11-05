# Module vars_regex

```
{
    "name": "",
    "source": "",
    "pattern": "",
    "overwrite: false,
}
```
## Description
____

An HTTP middleware module that evaluates regular expressions using placeholders or given text.

The results are also expressed as placeholders for use elsewhere in your route.

## Field List
____

name
  
A unique name for this expression that will be used to form the generated placeholder name(s)

source

A text string that is used as the evaluation source for the regex.  Strings enclosed in braces {} will be expanded as placeholders (if placeholders with the given name exist)

pattern

The golang regular expression (as described here: https://golang.org/s/re2syntax)

overwrite

Is a boolean flag that defaults to true.  If set it will overwrite placeholders values if they already exist.

## Examples

```
route {

    # Expected results:
    #   {"http.vars_regex.url_element.capture_group1": "https://"}
    #   {"http.vars_regex.url_element.capture_group2": "github.com"}
    #   {"http.vars_regex.url_element.capture_group3": ":443"}
    #   {"http.vars_regex.url_element.match1": "https://github.com:443/amalto"}

    vars_regex {
      name url_element
      source "https://github.com:443/amalto"
  	  pattern "(https?:\/\/)([^:^\/]*)(:\d*[^\/])?(.*[^\/])?"
    }

    # Expected results:
    #   {"http.vars_regex.named_url_element.scheme": "https://"}
    #   {"http.vars_regex.named_url_element.host": "github.com"}
    #   {"http.vars_regex.named_url_element.port": ":443"}
    #   {"http.vars_regex.named_url_element.path": "/amalto"}
    #   {"http.vars_regex.named_url_element.match1": "https://github.com:443/amalto"}

    vars_regex {
      name named_url_element
      source "https://github.com:443/amalto"
  	  pattern "(?P<scheme>https?:\/\/)(?P<host>[^:^\/]*)(?P<port>:\d*[^\/])?(?P<path>.*[^\/])?"
    }

    # Expected results:
    #   {"http.vars_regex.word1.match1": "{one}"}
    #   {"http.vars_regex.word1.match2": "{two}"}
    #   {"http.vars_regex.word1.match3": "{three}"}

    vars_regex {
        name word1
        source "{one} {two} {three}"
        pattern "{[\w.-]+}"
    }

    # Expected results:
    #   {"http.vars_regex.word2.capture_group1": "one"}
    #   {"http.vars_regex.word2.match1": "{one}"}
    #   {"http.vars_regex.word2.match2": "{two}"}
    #   {"http.vars_regex.word2.match3": "{three}"}

    vars_regex {
        name word2
        source "{one} {two} {three}"
        pattern "{([\w.-]+)}"
    }

    # Expected results:
    #   {"http.vars_regex.aport.match1": "443"}
    #   {"http.vars_regex.aport.capture_group1": ""}

    vars_regex {
        name aport
        source {http.vars_regex.url_element.capture_group3}
        pattern "(\d*)?"
    }

    # Expected results:
    #   {"http.vars_regex.aport.match1": "443"}   <-- not changed as overwrite == false
    #   {"http.vars_regex.aport.capture_group1": ""}

    vars_regex {
        name aport
        source ":8443"
        pattern "(\d*)?"
        overwrite false
    }

  }
```