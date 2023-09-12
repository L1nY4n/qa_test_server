<template>
    <!-- <v-chart ref="myChart" class="chart" :option="option" :update-options="update" autoresize /> -->
    <v-chart ref="myChart" class="chart" :option="option" ref= "mychart" autoresize />
</template>
  
<script lang="ts" setup>
  import { use } from 'echarts/core';
  import { CanvasRenderer } from 'echarts/renderers';
  import { PieChart } from 'echarts/charts';
  import { LineChart } from 'echarts/charts';
  import { UniversalTransition } from 'echarts/features';
  import {
    TitleComponent,
    TooltipComponent,
    LegendComponent,
  } from 'echarts/components';
  import VChart, { THEME_KEY,UPDATE_OPTIONS_KEY } from 'vue-echarts';
  import { ref, reactive,onMounted,provide } from 'vue';
  
  use([
    CanvasRenderer,
    PieChart,
    TitleComponent,
    TooltipComponent,
    LegendComponent,
  ]);
  
  provide(THEME_KEY, 'light');
  
  //provide(UPDATE_OPTIONS_KEY, 'randomData');

  function randomData() {
  now = new Date(+now + oneDay);
  value = value + Math.random() * 21 - 10;
  return {
    name: now.toString(),
    value: [
      [now.getFullYear(), now.getMonth() + 1, now.getDate()].join('/'),
      Math.round(value)
    ]
  };
}


let data :any= [];
let now = new Date(1997, 9, 3);
let oneDay = 24 * 3600 * 1000;
let value = Math.random() * 1000;

for (var i = 0; i < 1000; i++) {
  data.push(randomData());
}


  const option = (
    {
  title: {
    text: 'Dynamic Data & Time Axis'
  },
  tooltip: {
    trigger: 'axis',
    formatter: function (params:any) {
      params = params[0];
      var date = new Date(params.name);
      return (
        date.getDate() +
        '/' +
        (date.getMonth() + 1) +
        '/' +
        date.getFullYear() +
        ' : ' +
        params.value[1]
      );
    },
    axisPointer: {
      animation: false
    }
  },
  xAxis: {
    type: 'time',
    splitLine: {
      show: false
    }
  },
  yAxis: {
    type: 'value',
    boundaryGap: [0, '100%'],
    splitLine: {
      show: false
    }
  },
  series: [
    {
      name: 'Fake Data',
      type: 'line',
      showSymbol: false,
      data: data
    }
  ]
}

  );

 

  // const update=
  // {
  //   for (var i = 0; i < 5; i++) {
  //   data.shift();
  //   data.push(randomData());
  // }

   
  


  setInterval(function () {
  //  update();
       for (var i = 0; i < 5; i++) {
    data.shift();
    data.push(randomData());
       }
    
}, 1000);


  

 

//     setInterval(function () {
//   for (var i = 0; i < 5; i++) {
//     data.shift();
//     data.push(randomData());
//   }

// }, 1000);
  
// })
 

//   onMounted(() => {
//     setInterval(function () {
//   for (var i = 0; i < 5; i++) {
//     data.shift();
//     data.push(randomData());
//   }

// }, 1000);
  
// });

    
  



  </script>


  
  <style scoped>
  .chart {
     /* height: 100vh;  */
     height: 500px;
     width: 500px;
  }


  
  </style>