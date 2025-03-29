<template>
  <v-container>
    <v-row>
      <v-col lg="3" md="4" sm="6" cols="12" v-for="n in 10" :key="n">
        <v-card
          :to="{ name: 'recipe', params: { recipeId: n } }"
          title="Madre's Tamalles"
          hover
          subtitle="This is a subtitle"
          style="aspect-ratio: 1/1"
        >
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style></style>
