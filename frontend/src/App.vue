<template>
  <v-app>
    <v-app-bar
      app
      color="primary"
      dark
    >
      <div class="d-flex align-center">
        <v-app-bar-title
          dark
        >Image Share</v-app-bar-title>
      </div>
      <v-spacer></v-spacer>
    </v-app-bar>

    <v-main>
      <Images/>
    </v-main>
  </v-app>
</template>

<script>
import Images from './components/Images.vue';

export default {
  name: 'App',

  components: {
    Images
  },

  data: () => ({
    
  }),
};
</script>
