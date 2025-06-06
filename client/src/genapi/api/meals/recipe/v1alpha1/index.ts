// Code generated by protoc-gen-typescript-http. DO NOT EDIT.
/* eslint-disable camelcase */
// @ts-nocheck

// the main recipe object
export type Recipe = {
  // the name of the recipe
  //
  // Behaviors: IDENTIFIER
  name: string | undefined;
  // the title of the recipe
  //
  // Behaviors: REQUIRED
  title: string | undefined;
  // the description of the recipe
  //
  // Behaviors: OPTIONAL
  description: string | undefined;
  // the steps to make the recipe
  //
  // Behaviors: OPTIONAL
  directions: Recipe_Direction[] | undefined;
  // the ingredient groups in the recipe
  //
  // Behaviors: OPTIONAL
  ingredientGroups: Recipe_IngredientGroup[] | undefined;
  // image url
  //
  // Behaviors: OPTIONAL
  imageUri: string | undefined;
};

// the directions to make the recipe
export type Recipe_Direction = {
  // the title of the step
  //
  // Behaviors: OPTIONAL
  title: string | undefined;
  // the steps in the instruction
  //
  // Behaviors: REQUIRED
  steps: string[] | undefined;
};

export type Recipe_IngredientGroup = {
  // the name of the group
  //
  // Behaviors: OPTIONAL
  title: string | undefined;
  // the ingredients in the group
  //
  // Behaviors: REQUIRED
  ingredients: Recipe_Ingredient[] | undefined;
};

// an ingredient in a recipe
export type Recipe_Ingredient = {
  // the name of the ingredient
  //
  // Behaviors: REQUIRED
  title: string | undefined;
  // wheter the ingredient is optional
  //
  // Behaviors: OPTIONAL
  optional: boolean | undefined;
  // the quantity of the ingredient
  //
  // Behaviors: REQUIRED
  measurementAmount: number | undefined;
  // the type of measurement
  //
  // Behaviors: REQUIRED
  measurementType: Recipe_MeasurementType | undefined;
};

// the type of measurement
export type Recipe_MeasurementType =
  // the measurement is in cups
  | "MEASUREMENT_TYPE_UNSPECIFIED"
  // the measurement is in tablespoons
  | "MEASUREMENT_TYPE_TABLESPOON"
  // the measurement is in teaspoons
  | "MEASUREMENT_TYPE_TEASPOON"
  // the measurement is in ounces
  | "MEASUREMENT_TYPE_OUNCE"
  // the measurement is in pounds
  | "MEASUREMENT_TYPE_POUND"
  // the measurement is in grams
  | "MEASUREMENT_TYPE_GRAM"
  // the measurement is in milliliters
  | "MEASUREMENT_TYPE_MILLILITER"
  // the measurement is in liters
  | "MEASUREMENT_TYPE_LITER";
// the request to create a recipe
export type CreateRecipeRequest = {
  // the parent of the recipe
  //
  // Behaviors: OPTIONAL
  parent: string | undefined;
  // the recipe to create
  //
  // Behaviors: REQUIRED
  recipe: Recipe | undefined;
  // the id of the recipe
  //
  // Behaviors: REQUIRED
  recipeId: string | undefined;
};

// the request to list recipes
export type ListRecipesRequest = {
  // the parent of the recipe
  //
  // Behaviors: OPTIONAL
  parent: string | undefined;
  // returned page
  //
  // Behaviors: OPTIONAL
  pageSize: number | undefined;
  // used to specify the page token
  //
  // Behaviors: OPTIONAL
  pageToken: string | undefined;
  // used to specify the filter
  //
  // Behaviors: OPTIONAL
  filter: string | undefined;
};

// the response to list recipes
export type ListRecipesResponse = {
  // the recipes
  recipes: Recipe[] | undefined;
  // the next page token
  nextPageToken: string | undefined;
};

// the request to update a recipe
export type UpdateRecipeRequest = {
  // the recipe to update
  //
  // Behaviors: REQUIRED
  recipe: Recipe | undefined;
  // the fields to update
  //
  // Behaviors: OPTIONAL
  updateMask: wellKnownFieldMask | undefined;
};

// In JSON, a field mask is encoded as a single string where paths are
// separated by a comma. Fields name in each path are converted
// to/from lower-camel naming conventions.
// As an example, consider the following message declarations:
//
//     message Profile {
//       User user = 1;
//       Photo photo = 2;
//     }
//     message User {
//       string display_name = 1;
//       string address = 2;
//     }
//
// In proto a field mask for `Profile` may look as such:
//
//     mask {
//       paths: "user.display_name"
//       paths: "photo"
//     }
//
// In JSON, the same mask is represented as below:
//
//     {
//       mask: "user.displayName,photo"
//     }
type wellKnownFieldMask = string;

// the request to delete a recipe
export type DeleteRecipeRequest = {
  // the name of the recipe to delete
  //
  // Behaviors: REQUIRED
  name: string | undefined;
};

// the request to get a recipe
export type GetRecipeRequest = {
  // the name of the recipe to get
  //
  // Behaviors: REQUIRED
  name: string | undefined;
};

// the request to share a recipe
export type ShareRecipeRequest = {
  // the name of the recipe to share
  //
  // Behaviors: REQUIRED
  name: string | undefined;
  // the recipents of the recipe
  //
  // Behaviors: REQUIRED
  recipients: string[] | undefined;
  // the permission level given to the recipients
  //
  // Behaviors: REQUIRED
  permission: apitypes_PermissionLevel | undefined;
};

// the permission levels
export type apitypes_PermissionLevel =
  // the permission is not specified
  | "RESOURCE_PERMISSION_UNSPECIFIED"
  // the permission is read
  | "RESOURCE_PERMISSION_READ"
  // the permission is write
  | "RESOURCE_PERMISSION_WRITE";
// the response to share a recipe
export type ShareRecipeResponse = {
};

// the request to unshare a recipe
export type UnshareRecipeRequest = {
  // the name of the recipe to unshare
  //
  // Behaviors: REQUIRED
  name: string | undefined;
  // the recipients to remove from the recipe
  //
  // Behaviors: REQUIRED
  recipients: string[] | undefined;
};

// the response to unshare a recipe
export type UnshareRecipeResponse = {
};

// the recipe service
export interface RecipeService {
  // create a recipe
  CreateRecipe(request: CreateRecipeRequest): Promise<Recipe>;
  // list recipes
  ListRecipes(request: ListRecipesRequest): Promise<ListRecipesResponse>;
  // update a recipe
  UpdateRecipe(request: UpdateRecipeRequest): Promise<Recipe>;
  // delete` a recipe
  DeleteRecipe(request: DeleteRecipeRequest): Promise<Recipe>;
  // get a recipe
  GetRecipe(request: GetRecipeRequest): Promise<Recipe>;
  // share a recipe
  ShareRecipe(request: ShareRecipeRequest): Promise<ShareRecipeResponse>;
  // unshare a recipe
  UnshareRecipe(request: UnshareRecipeRequest): Promise<UnshareRecipeResponse>;
}

type RequestType = {
  path: string;
  method: string;
  body: string | null;
};

type RequestHandler = (request: RequestType, meta: { service: string, method: string }) => Promise<unknown>;

export function createRecipeServiceClient(
  handler: RequestHandler
): RecipeService {
  return {
    CreateRecipe(request) { // eslint-disable-line @typescript-eslint/no-unused-vars
      const path = `meals/v1alpha1/recipes`; // eslint-disable-line quotes
      const body = JSON.stringify(request?.recipe ?? {});
      const queryParams: string[] = [];
      if (request.parent) {
        queryParams.push(`parent=${encodeURIComponent(request.parent.toString())}`)
      }
      if (request.recipeId) {
        queryParams.push(`recipeId=${encodeURIComponent(request.recipeId.toString())}`)
      }
      let uri = path;
      if (queryParams.length > 0) {
        uri += `?${queryParams.join("&")}`
      }
      return handler({
        path: uri,
        method: "POST",
        body,
      }, {
        service: "RecipeService",
        method: "CreateRecipe",
      }) as Promise<Recipe>;
    },
    ListRecipes(request) { // eslint-disable-line @typescript-eslint/no-unused-vars
      const path = `meals/v1alpha1/recipes`; // eslint-disable-line quotes
      const body = null;
      const queryParams: string[] = [];
      if (request.parent) {
        queryParams.push(`parent=${encodeURIComponent(request.parent.toString())}`)
      }
      if (request.pageSize) {
        queryParams.push(`pageSize=${encodeURIComponent(request.pageSize.toString())}`)
      }
      if (request.pageToken) {
        queryParams.push(`pageToken=${encodeURIComponent(request.pageToken.toString())}`)
      }
      if (request.filter) {
        queryParams.push(`filter=${encodeURIComponent(request.filter.toString())}`)
      }
      let uri = path;
      if (queryParams.length > 0) {
        uri += `?${queryParams.join("&")}`
      }
      return handler({
        path: uri,
        method: "GET",
        body,
      }, {
        service: "RecipeService",
        method: "ListRecipes",
      }) as Promise<ListRecipesResponse>;
    },
    UpdateRecipe(request) { // eslint-disable-line @typescript-eslint/no-unused-vars
      if (!request.recipe?.name) {
        throw new Error("missing required field request.recipe.name");
      }
      const path = `meals/v1alpha1/${request.recipe.name}`; // eslint-disable-line quotes
      const body = JSON.stringify(request?.recipe ?? {});
      const queryParams: string[] = [];
      if (request.updateMask) {
        queryParams.push(`updateMask=${encodeURIComponent(request.updateMask.toString())}`)
      }
      let uri = path;
      if (queryParams.length > 0) {
        uri += `?${queryParams.join("&")}`
      }
      return handler({
        path: uri,
        method: "PATCH",
        body,
      }, {
        service: "RecipeService",
        method: "UpdateRecipe",
      }) as Promise<Recipe>;
    },
    DeleteRecipe(request) { // eslint-disable-line @typescript-eslint/no-unused-vars
      if (!request.name) {
        throw new Error("missing required field request.name");
      }
      const path = `meals/v1alpha1/${request.name}`; // eslint-disable-line quotes
      const body = null;
      const queryParams: string[] = [];
      let uri = path;
      if (queryParams.length > 0) {
        uri += `?${queryParams.join("&")}`
      }
      return handler({
        path: uri,
        method: "DELETE",
        body,
      }, {
        service: "RecipeService",
        method: "DeleteRecipe",
      }) as Promise<Recipe>;
    },
    GetRecipe(request) { // eslint-disable-line @typescript-eslint/no-unused-vars
      if (!request.name) {
        throw new Error("missing required field request.name");
      }
      const path = `meals/v1alpha1/${request.name}`; // eslint-disable-line quotes
      const body = null;
      const queryParams: string[] = [];
      let uri = path;
      if (queryParams.length > 0) {
        uri += `?${queryParams.join("&")}`
      }
      return handler({
        path: uri,
        method: "GET",
        body,
      }, {
        service: "RecipeService",
        method: "GetRecipe",
      }) as Promise<Recipe>;
    },
    ShareRecipe(request) { // eslint-disable-line @typescript-eslint/no-unused-vars
      if (!request.name) {
        throw new Error("missing required field request.name");
      }
      const path = `meals/v1alpha1/${request.name}:share`; // eslint-disable-line quotes
      const body = JSON.stringify(request);
      const queryParams: string[] = [];
      let uri = path;
      if (queryParams.length > 0) {
        uri += `?${queryParams.join("&")}`
      }
      return handler({
        path: uri,
        method: "POST",
        body,
      }, {
        service: "RecipeService",
        method: "ShareRecipe",
      }) as Promise<ShareRecipeResponse>;
    },
    UnshareRecipe(request) { // eslint-disable-line @typescript-eslint/no-unused-vars
      if (!request.name) {
        throw new Error("missing required field request.name");
      }
      const path = `meals/v1alpha1/${request.name}:unshare`; // eslint-disable-line quotes
      const body = JSON.stringify(request);
      const queryParams: string[] = [];
      let uri = path;
      if (queryParams.length > 0) {
        uri += `?${queryParams.join("&")}`
      }
      return handler({
        path: uri,
        method: "POST",
        body,
      }, {
        service: "RecipeService",
        method: "UnshareRecipe",
      }) as Promise<UnshareRecipeResponse>;
    },
  };
}

// @@protoc_insertion_point(typescript-http-eof)
