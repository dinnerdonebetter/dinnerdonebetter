package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertValidIngredientGroupCreationRequestInputToValidIngredientGroup(input *messages.ValidIngredientGroupCreationRequestInput) *messages.ValidIngredientGroup {
convertedmembers := make([]*messages.ValidIngredientGroupMember, 0, len(input.Members))
for _, item := range input.Members {
    convertedmembers = append(convertedmembers, ConvertValidIngredientGroupMemberCreationRequestInputToValidIngredientGroupMember(item))
}

output := &messages.ValidIngredientGroup{
    Members: convertedmembers,
    Name: input.Name,
    Slug: input.Slug,
    Description: input.Description,
}

return output
}
func ConvertValidIngredientGroupUpdateRequestInputToValidIngredientGroup(input *messages.ValidIngredientGroupUpdateRequestInput) *messages.ValidIngredientGroup {

output := &messages.ValidIngredientGroup{
    Name: input.Name,
    Slug: input.Slug,
    Description: input.Description,
}

return output
}
