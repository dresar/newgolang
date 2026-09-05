const { createApp, ref, computed, onMounted } = Vue
const { createRouter, createWebHashHistory } = VueRouter

const apiBase = "/api"

async function apiRequest(path, options = {}) {
  const res = await fetch(apiBase + path, {
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    ...options,
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    throw data
  }
  return data
}

const authState = {
  user: ref(null),
  loading: ref(false),
  initialized: ref(false),
}

async function fetchMe() {
  authState.loading.value = true
  try {
    const me = await apiRequest("/auth/me", { method: "GET" })
    authState.user.value = me
  } catch (e) {
    authState.user.value = null
  } finally {
    authState.loading.value = false
    authState.initialized.value = true
  }
}

async function login(username, password) {
  const data = await apiRequest("/auth/login", {
    method: "POST",
    body: JSON.stringify({ username, password }),
  })
  authState.user.value = data.user
}

async function logout() {
  await apiRequest("/auth/logout", { method: "POST" })
  authState.user.value = null
}

const LoginView = {
  setup() {
    const username = ref("admin")
    const password = ref("")
    const error = ref("")
    const loading = ref(false)
    const submit = async () => {
      error.value = ""
      loading.value = true
      try {
        await login(username.value, password.value)
        router.push("/dashboard")
      } catch (e) {
        error.value = e.error || "Login gagal"
      } finally {
        loading.value = false
      }
    }
    return { username, password, error, loading, submit }
  },
  template: `
    <div class="min-h-screen flex items-center justify-center">
      <div class="w-full max-w-md bg-slate-800 rounded-xl shadow-lg p-8">
        <h1 class="text-2xl font-bold mb-6 text-center">Login GoWa</h1>
        <div v-if="error" class="mb-4 text-red-400 text-sm">{{ error }}</div>
        <form @submit.prevent="submit" class="space-y-4">
          <div>
            <label class="block text-sm mb-1">Username</label>
            <input v-model="username" class="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 focus:outline-none" />
          </div>
          <div>
            <label class="block text-sm mb-1">Password</label>
            <input v-model="password" type="password" class="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 focus:outline-none" />
          </div>
          <button :disabled="loading" class="w-full py-2 rounded bg-emerald-500 hover:bg-emerald-600 text-white font-semibold">
            <span v-if="loading">Memproses...</span>
            <span v-else>Login</span>
          </button>
        </form>
      </div>
    </div>
  `,
}

const DashboardView = {
  setup() {
    const status = ref(null)
    const qr = ref("")
    const loadingStatus = ref(false)
    const loadingQR = ref(false)
    const error = ref("")
    const loadStatus = async () => {
      loadingStatus.value = true
      error.value = ""
      try {
        status.value = await apiRequest("/device/status", { method: "GET" })
      } catch (e) {
        error.value = e.error || "Gagal mengambil status"
      } finally {
        loadingStatus.value = false
      }
    }
    const loadQR = async () => {
      loadingQR.value = true
      error.value = ""
      try {
        const res = await apiRequest("/device/qr", { method: "GET" })
        qr.value = res.qr
      } catch (e) {
        qr.value = ""
        error.value = e.error || "QR belum tersedia"
      } finally {
        loadingQR.value = false
      }
    }
    const reconnect = async () => {
      error.value = ""
      try {
        await apiRequest("/device/reconnect", { method: "POST" })
        await loadStatus()
      } catch (e) {
        error.value = e.error || "Gagal reconnect"
      }
    }
    onMounted(() => {
      loadStatus()
      loadQR()
    })
    return { status, qr, error, loadingStatus, loadingQR, loadStatus, loadQR, reconnect }
  },
  template: `
    <div class="p-6 space-y-6">
      <h1 class="text-2xl font-bold">Dashboard</h1>
      <div v-if="error" class="text-red-400 text-sm">{{ error }}</div>
      <div class="grid md:grid-cols-2 gap-6">
        <div class="bg-slate-800 rounded-xl p-4 space-y-3">
          <div class="flex justify-between items-center">
            <h2 class="font-semibold">Status Device</h2>
            <button @click="loadStatus" class="text-xs px-2 py-1 rounded bg-slate-700 hover:bg-slate-600">
              Refresh
            </button>
          </div>
          <div v-if="loadingStatus">Memuat...</div>
          <div v-else-if="status">
            <div class="flex items-center space-x-2">
              <span class="w-2 h-2 rounded-full" :class="status.connected ? 'bg-emerald-400' : 'bg-red-400'"></span>
              <span>{{ status.connected ? 'Terhubung' : 'Tidak terhubung' }}</span>
            </div>
            <div class="text-sm text-slate-300 mt-2">
              <div>Event: {{ status.last_event }}</div>
              <div>Nomor: {{ status.phone_number || '-' }}</div>
            </div>
            <button @click="reconnect" class="mt-3 px-3 py-1 rounded bg-amber-500 hover:bg-amber-600 text-sm">
              Reconnect
            </button>
          </div>
        </div>
        <div class="bg-slate-800 rounded-xl p-4 space-y-3">
          <div class="flex justify-between items-center">
            <h2 class="font-semibold">QR Code Login</h2>
            <button @click="loadQR" class="text-xs px-2 py-1 rounded bg-slate-700 hover:bg-slate-600">
              Refresh QR
            </button>
          </div>
          <div v-if="loadingQR">Memuat QR...</div>
          <div v-else-if="qr" class="space-y-2">
            <div class="text-xs break-all bg-slate-900 rounded p-2 border border-slate-700">
              {{ qr }}
            </div>
            <div class="text-xs text-slate-400">
              Scan teks QR ini menggunakan kamera lain atau generate ke gambar dengan tool QR.
            </div>
          </div>
          <div v-else class="text-sm text-slate-400">
            QR belum tersedia atau sudah kadaluwarsa.
          </div>
        </div>
      </div>
    </div>
  `,
}

const SendMessageView = {
  setup() {
    const to = ref("")
    const message = ref("")
    const sending = ref(false)
    const error = ref("")
    const success = ref("")
    const submit = async () => {
      error.value = ""
      success.value = ""
      sending.value = true
      try {
        await apiRequest("/messages/send", {
          method: "POST",
          body: JSON.stringify({ to: to.value, message: message.value }),
        })
        success.value = "Pesan terkirim"
      } catch (e) {
        error.value = e.error || "Gagal mengirim pesan"
      } finally {
        sending.value = false
      }
    }
    return { to, message, sending, error, success, submit }
  },
  template: `
    <div class="p-6 space-y-4">
      <h1 class="text-2xl font-bold">Kirim Pesan</h1>
      <div v-if="error" class="text-red-400 text-sm">{{ error }}</div>
      <div v-if="success" class="text-emerald-400 text-sm">{{ success }}</div>
      <form @submit.prevent="submit" class="space-y-4 max-w-xl">
        <div>
          <label class="block text-sm mb-1">Nomor Tujuan (tanpa +, contoh 62812xxxx)</label>
          <input v-model="to" class="w-full px-3 py-2 rounded bg-slate-800 border border-slate-700 focus:outline-none" />
        </div>
        <div>
          <label class="block text-sm mb-1">Pesan</label>
          <textarea v-model="message" rows="4" class="w-full px-3 py-2 rounded bg-slate-800 border border-slate-700 focus:outline-none"></textarea>
        </div>
        <button :disabled="sending" class="px-4 py-2 rounded bg-emerald-500 hover:bg-emerald-600 text-sm font-semibold">
          <span v-if="sending">Mengirim...</span>
          <span v-else>Kirim</span>
        </button>
      </form>
    </div>
  `,
}

const AutoReplyView = {
  setup() {
    const rules = ref([])
    const loading = ref(false)
    const error = ref("")
    const form = {
      keyword: ref(""),
      reply_text: ref(""),
      match_type: ref("contains"),
      is_active: ref(true),
    }
    const loadRules = async () => {
      loading.value = true
      error.value = ""
      try {
        rules.value = await apiRequest("/autoreply/", { method: "GET" })
      } catch (e) {
        error.value = e.error || "Gagal mengambil data"
      } finally {
        loading.value = false
      }
    }
    const createRule = async () => {
      error.value = ""
      try {
        await apiRequest("/autoreply/", {
          method: "POST",
          body: JSON.stringify({
            keyword: form.keyword.value,
            reply_text: form.reply_text.value,
            match_type: form.match_type.value,
            is_active: form.is_active.value,
          }),
        })
        form.keyword.value = ""
        form.reply_text.value = ""
        await loadRules()
      } catch (e) {
        error.value = e.error || "Gagal menyimpan"
      }
    }
    const toggleActive = async (rule) => {
      try {
        await apiRequest("/autoreply/" + rule.id, {
          method: "PUT",
          body: JSON.stringify({ is_active: !rule.is_active }),
        })
        await loadRules()
      } catch (e) {
        error.value = e.error || "Gagal update"
      }
    }
    const removeRule = async (rule) => {
      try {
        await apiRequest("/autoreply/" + rule.id, { method: "DELETE" })
        await loadRules()
      } catch (e) {
        error.value = e.error || "Gagal hapus"
      }
    }
    onMounted(loadRules)
    return { rules, loading, error, form, loadRules, createRule, toggleActive, removeRule }
  },
  template: `
    <div class="p-6 space-y-6">
      <h1 class="text-2xl font-bold">Auto Reply</h1>
      <div v-if="error" class="text-red-400 text-sm">{{ error }}</div>
      <div class="grid md:grid-cols-2 gap-6">
        <div class="bg-slate-800 rounded-xl p-4 space-y-3">
          <h2 class="font-semibold text-sm">Tambah Rule</h2>
          <div class="space-y-3">
            <div>
              <label class="block text-xs mb-1">Keyword</label>
              <input v-model="form.keyword" class="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-sm" />
            </div>
            <div>
              <label class="block text-xs mb-1">Balasan</label>
              <textarea v-model="form.reply_text" rows="3" class="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-sm"></textarea>
            </div>
            <div class="flex items-center space-x-3 text-xs">
              <select v-model="form.match_type" class="px-2 py-1 rounded bg-slate-900 border border-slate-700">
                <option value="contains">Mengandung</option>
                <option value="equals">Sama persis</option>
              </select>
              <label class="inline-flex items-center space-x-2">
                <input type="checkbox" v-model="form.is_active" class="rounded border-slate-700 bg-slate-900" />
                <span>Aktif</span>
              </label>
            </div>
            <button @click="createRule" class="px-3 py-1 rounded bg-emerald-500 hover:bg-emerald-600 text-xs font-semibold">
              Simpan
            </button>
          </div>
        </div>
        <div class="bg-slate-800 rounded-xl p-4 space-y-3">
          <div class="flex justify-between items-center">
            <h2 class="font-semibold text-sm">Daftar Rule</h2>
            <button @click="loadRules" class="text-xs px-2 py-1 rounded bg-slate-700 hover:bg-slate-600">
              Refresh
            </button>
          </div>
          <div v-if="loading" class="text-sm">Memuat...</div>
          <div v-else-if="rules.length === 0" class="text-sm text-slate-400">
            Belum ada rule.
          </div>
          <div v-else class="space-y-2 max-h-80 overflow-y-auto">
            <div v-for="r in rules" :key="r.id" class="border border-slate-700 rounded p-2 text-xs flex justify-between items-start">
              <div>
                <div class="font-semibold">{{ r.keyword }}</div>
                <div class="text-slate-300 whitespace-pre-wrap">{{ r.reply_text }}</div>
                <div class="mt-1 text-[10px] text-slate-400">
                  Mode: {{ r.match_type }} | Status: {{ r.is_active ? 'Aktif' : 'Nonaktif' }}
                </div>
              </div>
              <div class="space-y-1 ml-2 flex flex-col items-end">
                <button @click="toggleActive(r)" class="px-2 py-1 rounded bg-slate-700 hover:bg-slate-600">
                  {{ r.is_active ? 'Nonaktif' : 'Aktifkan' }}
                </button>
                <button @click="removeRule(r)" class="px-2 py-1 rounded bg-red-600 hover:bg-red-700">
                  Hapus
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  `,
}

const Layout = {
  setup() {
    const user = computed(() => authState.user.value)
    const loggingOut = ref(false)
    const doLogout = async () => {
      loggingOut.value = true
      try {
        await logout()
        router.push("/login")
      } finally {
        loggingOut.value = false
      }
    }
    return { user, loggingOut, doLogout }
  },
  template: `
    <div class="min-h-screen flex">
      <aside class="w-56 bg-slate-950 border-r border-slate-800 p-4 space-y-4 hidden md:block">
        <div class="text-lg font-bold">GoWa Vue Lite</div>
        <nav class="space-y-2 text-sm">
          <RouterLink to="/dashboard" class="block px-3 py-2 rounded hover:bg-slate-800" active-class="bg-slate-800">
            Dashboard
          </RouterLink>
          <RouterLink to="/send" class="block px-3 py-2 rounded hover:bg-slate-800" active-class="bg-slate-800">
            Kirim Pesan
          </RouterLink>
          <RouterLink to="/autoreply" class="block px-3 py-2 rounded hover:bg-slate-800" active-class="bg-slate-800">
            Auto Reply
          </RouterLink>
        </nav>
        <div class="mt-auto text-xs text-slate-400 space-y-2">
          <div v-if="user">Login sebagai {{ user.username }}</div>
          <button @click="doLogout" :disabled="loggingOut" class="px-3 py-1 rounded bg-slate-800 hover:bg-slate-700 w-full text-left">
            Logout
          </button>
        </div>
      </aside>
      <div class="flex-1 flex flex-col">
        <header class="md:hidden flex items-center justify-between px-4 py-3 border-b border-slate-800 bg-slate-950">
          <div class="font-semibold">GoWa Vue Lite</div>
          <div class="flex items-center space-x-3 text-xs">
            <span v-if="user">Hi, {{ user.username }}</span>
            <button @click="doLogout" :disabled="loggingOut" class="px-3 py-1 rounded bg-slate-800 hover:bg-slate-700">
              Logout
            </button>
          </div>
        </header>
        <main class="flex-1">
          <RouterView />
        </main>
      </div>
    </div>
  `,
}

const routesDef = [
  { path: "/login", component: LoginView },
  { path: "/dashboard", component: DashboardView },
  { path: "/send", component: SendMessageView },
  { path: "/autoreply", component: AutoReplyView },
  { path: "/:pathMatch(.*)*", redirect: "/dashboard" },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes: routesDef,
})

router.beforeEach(async (to, from, next) => {
  if (!authState.initialized.value) {
    await fetchMe()
  }
  const isLoggedIn = !!authState.user.value
  if (to.path === "/login" && isLoggedIn) {
    next("/dashboard")
  } else if (to.path !== "/login" && !isLoggedIn) {
    next("/login")
  } else {
    next()
  }
})

const app = createApp(Layout)
app.use(router)
app.mount("#app")

