/*
Package apiclient provides an HTTP client that can communicate with and interpret the responses of an instance of the service.
*/
package apiclient

/* DELETEME

find: var (\w+) (\*?)types\.APIResponse\[(\[\])?\*types\.(\w+)\]\n\tif err = c.fetchAndUnmarshal\(ctx\, req\, &(\w+)\)\; err \!= nil \{\n\t\treturn nil, observability.Prepare(AndLog)?Error\(err, logger, span, "([\w\s\%]+)"(\, \w+)?\)\n\t\}

replace: var $1 $2types.APIResponse[$3*types.$4]
if err = c.unmarshalBody(ctx, res, &$5); err != nil {
        return nil, observability.Prepare$6Error(err, logger, span, "$7")
}

*/
