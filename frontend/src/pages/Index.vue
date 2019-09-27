<template>
  <q-page class="q-pa-lg">
    <q-form class="row no-wrap fullwidth items-baseline" @submit="onSubmit">
      <q-input class="col" label="From" v-model="form.from" mask="date" :rules="['date']">
        <template v-slot:append>
          <q-icon name="event" class="cursor-pointer">
            <q-popup-proxy ref="qDateProxy" transition-show="scale" transition-hide="scale">
              <q-date v-model="form.from" @input="() => $refs.qDateProxy.hide()" />
            </q-popup-proxy>
          </q-icon>
        </template>
      </q-input>

      <q-input class="col" label="To" v-model="form.to" mask="date" :rules="['date']">
        <template v-slot:append>
          <q-icon name="event" class="cursor-pointer">
            <q-popup-proxy ref="qDateProxy" transition-show="scale" transition-hide="scale">
              <q-date v-model="form.to" @input="() => $refs.qDateProxy.hide()" />
            </q-popup-proxy>
          </q-icon>
        </template>
      </q-input>

      <q-select
        class="col"
        label="Machines"
        multiple
        v-model="form.machinesSelected"
        :options="form.machines"
      />

      <q-select class="col" label="Sets" multiple v-model="form.setsSelected" :options="form.sets" />

      <q-btn class="col" label="Submit" small flat color="primary" type="submit" />
    </q-form>

    <q-table :data="qData.mainTable" :columns="table.columns" row-key="name" flat />
  </q-page>
</template>

<script>
import { date } from "quasar";
import axios from "axios";

const host = axios.create({
  baseURL: "http://localhost:3002/"
});

export default {
  name: "PageIndex",

  data() {
    return {
      form: {
        from: "2015-10-01",
        to: date.formatDate(new Date(), "YYYY-MM-DD"),
        sets: [],
        machines: [],
        setsSelected: ["all"],
        machinesSelected: ["all"]
      },
      qData: [],
      table: {
        columns: [
          {
            name: "label",
            align: "left",
            label: "Draw",
            field: "label"
          },
          {
            name: "b1",
            align: "left",
            label: "Ball 1",
            field: row => row.num[0]
          },
          {
            name: "b2",
            align: "left",
            label: "Ball 2",
            field: row => row.num[1]
          },
          {
            name: "b3",
            align: "left",
            label: "Ball 3",
            field: row => row.num[2]
          },
          {
            name: "b4",
            align: "left",
            label: "Ball 4",
            field: row => row.num[3]
          },
          {
            name: "b5",
            align: "left",
            label: "Ball 5",
            field: row => row.num[4]
          },
          {
            name: "b6",
            align: "left",
            label: "Ball 6",
            field: row => row.num[5]
          },
          {
            name: "bonus",
            align: "left",
            label: "Bonus",
            field: row => row.num[6]
          }
        ]
      }
    };
  },

  computed: {
    params() {
      return {
        start: new Date(this.form.from),
        end: new Date(this.form.to),
        sets: this.setsList,
        machines: this.machinesList
      };
    },
    setsList() {
      const out = [];
      for (let item of this.form.setsSelected) {
        if (item === "all") {
          return [];
        }

        const i = parseInt(item);
        if (!isNaN(i)) {
          out.push(i);
        }
      }
      return out;
    },
    machinesList() {
      const out = [];
      for (let item of this.form.machinesSelected) {
        if (item === "all") {
          return [];
        }
        out.push(item);
      }
      return out;
    }
  },

  created() {
    this.updateOpts();
  },

  methods: {
    onSubmit() {
      host
        .post("/query", this.params)
        .then(res => {
          console.log(res.data);
          this.qData = res.data;
        })
        .catch(err => alert(err));
    },
    updateOpts() {
      host
        .post("/sets", this.params)
        .then(res => {
          this.form.sets = res.data;
          this.form.sets.unshift("all");
        })
        .catch(err => alert(err));

      host
        .post("/machines", this.params)
        .then(res => {
          this.form.machines = res.data;
          this.form.machines.unshift("all");
        })
        .catch(err => alert(err));
    }
  }
};
</script>
