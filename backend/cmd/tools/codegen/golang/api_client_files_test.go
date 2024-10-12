package golang

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_apiClientFunction_Render(T *testing.T) {
	T.Parallel()

	T.Run("query filtered", func(t *testing.T) {
		t.Parallel()

		x := APIClientFunction{
			Name:          "makeTheSandwich",
			QueryFiltered: true,
			PathTemplate:  "/api/v1/stuff/{thing}/dedupe",
			Params: []functionParam{
				{
					Name:         "thing",
					Type:         "string",
					DefaultValue: `"baloney"`,
				},
				{
					Name: "stuff",
					Type: "bool",
				},
			},
			ResponseType: functionResponseType{
				GenericContainer: "APIResponse",
				TypeName:         "Sandwich",
			},
		}

		expected := `export async function makeTheSandwich(
  client: Axios,
  filter: QueryFilter = QueryFilter.Default(),
  thing: string = "baloney",
	stuff: bool,
	): Promise< APIResponse <  Sandwich  >  > {
  return new Promise(async function (resolve, reject) {
    const response = await client.get< APIResponse < Sandwich[]  >  >(` + "`" + `/api/v1/stuff/{thing}/dedupe` + "`" + `, {
      params: filter.asRecord(),
    });

    if (response.data.error) {
      reject(new Error(response.data.error.message));
    }

    const result = new QueryFilteredResult<Sandwich>({
      data: response.data.data,
      totalCount: response.data.pagination?.totalCount,
      page: response.data.pagination?.page,
      limit: response.data.pagination?.limit,
    });

    resolve(result);
  });
}`
		actual, _, err := x.Render()
		require.NoError(t, err)

		assert.Equal(t, expected, actual)
	})
}
