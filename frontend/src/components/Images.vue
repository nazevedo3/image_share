<template>
  <v-data-table
    :headers="headers"
    :items="images"
    sort-by="title"
    class="elevation-1"
    :search="search"
    show-expand
  >
  <template v-slot:[`item.path`]="{ item }"><v-img :src="item.path" style="height: 50px; width:50px"/></template>
  <template v-slot:[`expanded-item`]="{ headers, item }">
      <td :colspan="headers.length">
          <v-row>
              <v-col cols="12" md="4">
                <v-img :src="item.path" style="height: 400px; width:400px"></v-img>
              </v-col>
                <v-card-text class="font-italic">
                  <a :href="item.path" target="_blank">{{item.path}}</a>
                </v-card-text>
            </v-row>
        </td>
  </template>
  <template v-slot:[`item.tags`]="{ item }">
      <span v-for="(tag, index) in item.tags" :key="index">
          <v-chip
            class="ma-2"
            color="light blue"
            text-color="white"
            small
            >
            {{tag}}
          </v-chip>
      </span>


  </template>
    <template v-slot:top>
      <v-toolbar
        flat
      >
        <v-toolbar-title>Your Images</v-toolbar-title>
        <v-divider
          class="mx-4"
          inset
          vertical
        ></v-divider>
        <v-spacer></v-spacer>

        <v-card style="position:relative; left:-200px">
      <v-spacer></v-spacer>
      <v-text-field
        v-model="search"
        append-icon="mdi-magnify"
        label="Search Tags"
        single-line
        hide-details
      ></v-text-field>
      </v-card>



        <v-dialog
          v-model="dialog"
          max-width="500px"
        >
          <template v-slot:activator="{ on, attrs }">
            <v-btn
              color="primary"
              dark
              class="mb-2"
              v-bind="attrs"
              v-on="on"
            >
              Upload Image
            </v-btn>
          </template>
    
          <v-card>
            <v-card-title>
              <span class="text-h5">{{ formTitle }}</span>
            </v-card-title>

            <v-card-text>
              <v-container>
                <v-form ref="form" class="mx-2" lazy-validation>
                <v-row>
                  <v-col
                    cols="12"
                    sm="6"
                    md="4"
                  >
                    <v-text-field
                      v-model="editedItem.title"
                      label="Title"
                      required
                      :rules="titleRules"
                    hint="15 character max"
                    persistent-hint
                    ></v-text-field>
                  </v-col>
                  <v-col
                    cols="12"
                    sm="6"
                    md="4"
                  >
                    <v-text-field
                      v-model="editedItem.description"
                      label="Description"
                      required
                      :rules="descriptionRules"
                      hint="20 character max"
                      persistent-hint
                    ></v-text-field>
                  </v-col>
                  <v-col
                    cols="12"
                    sm="6"
                    md="4"
                  >
                    <v-combobox
                      multiple
                      v-model="editedItem.tags"
                      small
                      label="Tags"
                      append-icon
                      chips
                      deletable-chips
                      class="tag-input"
                      :search-input.sync="tagSearch"
                      @keyup.tab="updateTags"
                      @paste="updateTags"
                      required
                      :rules="tagRules"
                    ></v-combobox>
                  </v-col>
                  <v-col
                    cols="12"
                    sm="10"
                    md="10"
                  >
                    <v-file-input
                      @change="fileSelected"
                      :rules="rules"
                      accept="image/png, image/jpeg, image/bmp, image/gif"
                      placeholder="Select Image"
                      prepend-icon="mdi-camera"
                      label="Image"
                    ></v-file-input>
                  </v-col>
                </v-row>
                </v-form>
              </v-container>
            </v-card-text>

            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn
                color="blue darken-1"
                text
                @click="close"
              >
                Cancel
              </v-btn>
              <v-btn
                color="blue darken-1"
                text
                @click="save"
              >
                Save
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-toolbar>
    </template>
    <template v-slot:[`item.actions`]="{ item }">
      <v-icon
        small
        @click="deleteItem(item)"
      >
        mdi-delete
      </v-icon>
    </template>
  </v-data-table>
 </template>


<script>
import axios from '@/axios';

  export default {
    data: () => ({
      dialog: false,
      rules: [
        value => !value || value.size < 2000000 || 'Avatar size should be less than 2 MB!',
      ],
      select: [],
      tagSearch: '',
      snackbar: {
        appear: false,
        text: '',
        color: ''
      },
      headers: [
        {
          text: 'Image',
          align: 'start',
          sortable: false,
          value: 'path',
        },
        { text: 'Title', value: 'title' },
        { text: 'Description', value: 'description' },
        { text: 'Tags', value: 'tags' },
        { text: 'Delete', value: 'actions', sortable: false },
      ],
      selectedFile: null,
      images: [],
      search: '',
      titleRules: [ v => !!v || 'Title is required', 
                   v => (v && v.length <= 10) || 'Title must be less than 10 characters'],
      descriptionRules: [ v => !!v || 'Description is required',
                          v => (v && v.length <= 20) || 'Description must be less than 20 characters'],
      tagRules: [ v => (v && v.length <= 3) || 'Tags max is 3'],
      editedItem: {
        title: '',
        description: '',
        tags: []
      },
      defaultItem: {
        title: '',
        description: '',
        tags: []
      },
    }),

    computed: {
      formTitle () {
        return 'New Image'
      },
    },

    watch: {
      dialog (val) {
        val || this.close()
      },
    },

    created () {
      this.getAllImages()
    },

    methods: {
      fileSelected(file){
        this.selectedFile = file
      },
        async getAllImages(){
            await axios.get('api/images')
                .then(response => {
                    this.images = response.data
                    console.log("Got all Images: ", this.images)
                })
            .catch(e => {
                console.log(e)
            })
        },
      async deleteItem (item) {
            console.log("IMAGE TO BE DELETED", item)
            await axios.delete('api/images/'+item.id, {
                data: item,
                Headers: {
                    'Content-Type': 'multipart/form-data'
                }
            })
            .then(response => {
                console.log(response)
                this.getAllImages()

            })
            .catch(e => {
                this.errors.push(e)
            })
        },
      close () {
        this.dialog = false
        this.$refs.form.reset()
      

      },
      async save () {
        if (this.$refs.form.validate()){
        const fd = new FormData()
        fd.append('myFile', this.selectedFile, this.selectedFile.name)
        const data = JSON.stringify({title: this.editedItem.title, description: this.editedItem.description, tags: this.editedItem.tags})
        fd.append('data', data)
        await axios.post('api/images', fd, {
          Headers: {
            'Content-Type': 'multipart/form-data'
            }
            })
            .then(response => {
              console.log("Success: ", response)
              this.close()
              this.getAllImages()
              
            })
            .catch(e => {
                console.log("There was an error: ", e)
            })
        }
        this.$refs.form.reset()
      },
      updateTags(){
          this.$nextTick(() =>{
              this.select.push(...this.editedItem.tags.split(","));
              this.$nextTicket(() => {
                  this.editedItem.tags = "";
              })
          })
      }
    },
  }
</script>