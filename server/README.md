
## Architecture Overview
```plantuml
left to right direction
package "authentication" {
  package "api" {
    class oauth2 {
      create/login user then create/save
      a token on the server for the user. 
      Return a "token key" that is used to
      retrieve the token.
    }
    class GetToken {
      Get a token from the server and return it.
      This will delete it from the server. The token
      key can only be used once.
      ==
      token_key
    }
    class check {
      return the user "name" of 
      the authenticated user
      from the token on the
      Authorization header.
    }
  }
  package "resources" {
    class rToken {
      name STRING
      email STRING
    }
  }
  GetToken --> rToken
  oauth2 --> rToken
}

package "users" {
  package "api" {
    class GetUser {
      user_name STRING
    }
  }
  package "resources" {
    class rUser {
      name STRING
      email STRING
    }
  }
  package "database" {
    class dUser {
      user_id BIGINT
      email VARCHAR
      google_id VARCHAR
      facebook_id VARCHAR
      amazon_id VARCHAR
      --
      create_time TIMESTAMP
      update_time TIMESTAMP
      delete_time TIMESTAMP
    }
    
  }
  GetUser -- rUser
  rUser -- dUser
}

package "meals" {
  package "api" {
    class GetRecipe {
      Get a recipe. The user needs to have
      read access to the recipe.
      ==
      name STRING
    }
    class ListRecipes {
      ListR a recipes. The user needs to have
      read access to the recipes.
      ==
      parent STRING
    }
    class CreateRecipe {
      Create a recipe.
      ==
      parent STRING
      title STRING
      description STRING
      instructions Instructions
      ingredients Ingredients
    }
    class UpdateRecipe {
      Update a recipe. The user needs to have
      write access to the recipe.
      ==
      name STRING
      title STRING
      description STRING
      instructions Instructions
      ingredients Instructions
    }
    class DeleteRecipe {
      Delete a recipe. The user needs to have
      write access to the recipe.
      ==
      name STRING
    }
    class ShareRecipe {
      Share a recipe. The user needs to have
      write access to the recipe.
      ==
      name STRING
      users repeated STRING
      permission_level ENUM
    }
  }
  package "resources" {
    class rRecipe {
      name STRING
      title STRING
      description STRING
      instructions Instructions
      ingredients Instructions
    }
  }
  package "database" {

    class dRecipe {
      recipe_id BIGINT
      title VARCHAR
      description VARCHAR
      instructions BJSON
      --
      create_time TIMESTAMP
      update_time TIMESTAMP
      delete_time TIMESTAMP
    }
    class dRecipeIngredient {
      recipe_ingredient_id BIGINT
      recipe_id BIGINT
      ingredient_id BIGINT
      measurement BJSON
      optional BOOLEAN
      --
      create_time TIMESTAMP
      update_time TIMESTAMP
      delete_time TIMESTAMP
    }
    class dIngredient {
      ingredient_id BIGINT
      title VARCHAR
      --
      create_time TIMESTAMP
      update_time TIMESTAMP
      delete_time TIMESTAMP
    }
    class dRecipeUser {
      user_recipe_id BIGINT
      recipe_id BIGINT
      user_id BIGINT
      permission_level INT
      --
      create_time TIMESTAMP
      update_time TIMESTAMP
      delete_time TIMESTAMP
    }
    dRecipeIngredient - dRecipe
    dRecipe - dRecipeUser
    dIngredient - dRecipeIngredient

    GetRecipe --> rRecipe
    ListRecipes --> rRecipe
    CreateRecipe --> rRecipe
    UpdateRecipe --> rRecipe
    DeleteRecipe --> rRecipe
    ShareRecipe --> rRecipe
  }
  rRecipe -- dRecipe
  rRecipe -- dRecipeIngredient
  rRecipe -- dRecipeUser
  rRecipe -- dIngredient
}
dUser - dRecipeUser
oauth2 --> rUser
```