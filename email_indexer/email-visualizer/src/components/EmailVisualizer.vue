<template>
  <div class="flex h-4/5" ref="container">
    <div class="w-1/2 bg-gray-100 max-h-screen overflow-y-auto" style="overflow-y: auto">
      <table class="w-full">
        <thead>
          <tr>
            <th class="px-4 py-2">Subject</th>
            <th class="px-4 py-2">From</th>
            <th class="px-4 py-2">To</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="email in emailList" :key="email.id" @click="showEmailDetails(email)">
            <td class="border px-4 py-2">{{ email._source.subject }}</td>
            <td class="border px-4 py-2">{{ email._source.from }}</td>
            <td class="border px-4 py-2">{{ email._source.to }}</td>
          </tr>
        </tbody>
        <Observer @intersect="intersected" />
      </table>
    </div>

    <div class="w-1/2 bg-gray-200 max-h-screen overflow-y-auto">
      <div class="p-6">
        <h2 class="text-xl font-semibold mb-4">Selected Email</h2>
        <div class="border border-gray-300 p-4 rounded-lg">
          <template v-if="selectedEmail?._source">
            <p v-if="selectedEmail._source.from">
              <strong>From:</strong> {{ selectedEmail._source.from }}
            </p>
            <p v-if="selectedEmail._source.date">
              <strong>Date:</strong> {{ selectedEmail._source.date }}
            </p>
            <p v-if="selectedEmail._source.to">
              <strong>To:</strong> {{ selectedEmail._source.to }}
            </p>
            <p v-if="selectedEmail._source.bcc">
              <strong>BCC:</strong> {{ selectedEmail._source.bcc }}
            </p>
            <p v-if="selectedEmail._source.cc">
              <strong>CC:</strong> {{ selectedEmail._source.cc }}
            </p>
            <p v-if="selectedEmail._source.x_bcc">
              <strong>X_BCC:</strong> {{ selectedEmail._source.x_bcc }}
            </p>
            <p v-if="selectedEmail._source.x_cc">
              <strong>X_CC:</strong> {{ selectedEmail._source.x_cc }}
            </p>
            <p v-if="selectedEmail._source.x_from">
              <strong>X_From:</strong> {{ selectedEmail._source.x_from }}
            </p>
            <p v-if="selectedEmail._source.x_to">
              <strong>X_To:</strong> {{ selectedEmail._source.x_to }}
            </p>
            <p v-if="selectedEmail._source.subject">
              <strong>Subject:</strong> {{ selectedEmail._source.subject }}
            </p>
            <p v-if="selectedEmail._source.content">
              <strong>Content:</strong> {{ selectedEmail._source.content }}
            </p>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, watch } from 'vue'
import Observer from './Observer.vue'

export default {
  props: ['emails'],
  data: () => ({
    observer: null
  }),
  components: {
    Observer
  },
  methods: {},
  setup(props, ctx) {
    const emailList = ref([])
    const selectedEmail = ref(null)
    const container = ref(null)
    watch(
      () => props.emails,
      (newValue) => {
        emailList.value = newValue
      }
    )

    const showEmailDetails = (email) => {
      selectedEmail.value = email
    }

    const intersected = () => {
      ctx.emit('intersected')
    }
    return {
      emailList,
      selectedEmail,
      showEmailDetails,
      container,
      Observer,
      intersected
    }
  }
}
</script>
