<script setup lang="ts">
import { useI18n } from "vue-i18n";

import { ref, onMounted, onBeforeUnmount } from 'vue';

import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime';
import { SelectFile } from '../../wailsjs/go/main/App';
import { RunSeqAnalysis } from '../../wailsjs/go/main/App';


const { t } = useI18n();

const filePath = ref('');
const workdir = ref('')
const outputPrefix = ref('')
const noPolyA = ref(false)
const lessMem = ref(false)
const rowsLimit = ref(100000)

const messages = ref<string[]>([])

const dirname = (filePath: string) => {
  // 正则表达式，用于匹配文件路径分隔符（Windows和Unix/Linux）
  const separator = /[\\\/]/;
  // 移除文件路径末尾的目录分隔符（如果存在）
  const normalizedPath = filePath.replace(/[\\\/]*$/, '');
  // 分割路径
  const segments = normalizedPath.split(separator);
  // 移除路径的最后一个部分（文件名）
  segments.pop(); // 如果数组为空（即文件路径为根目录），pop操作将不会改变数组

  // 如果结果是空数组，表示我们已经处于根路径，返回原路径的根部分
  if (segments.length === 0) {
    return filePath.startsWith('/') ? '/' : '.';
  }

  // 重新组合剩余的路径部分
  return segments.join('/');
}

const basename = (filePath: string) => {
  // 确保文件路径是字符串
  if (typeof filePath !== 'string') {
    throw new TypeError('Path must be a string');
  }

  // 正则表达式，用于匹配文件路径分隔符（Windows和Unix/Linux）
  const separator = /[\\\/]/;
  // 移除文件路径末尾的目录分隔符（如果存在）
  const normalizedPath = filePath.replace(/[\\\/]*$/, '');
  // 分割路径
  const segments = normalizedPath.split(separator);
  // 提取出基本文件名（basename）
  const base = segments.pop() || ''; // 如果数组为空，则基本文件名是空字符串

  return base;
}

const openFileDialog = async () => {
  try {
    const options = {
      Title: '选择文件',
      Filters: '',
    };
    const result = await SelectFile(options.Title);
    filePath.value = result;
    workdir.value = dirname(filePath.value)
    outputPrefix.value = workdir.value + "/" + basename(workdir.value)
    // 这里可以添加更多处理文件路径的逻辑
  } catch (error) {
    console.error('Error selecting file:', error);
  }
};


const runSeqAnalysis = async () => {
  messages.value = [];

  // parse args
  const args = ["-zip"];
  if (filePath.value == "") {
    alert("Please select a file")
    return
  } else {
    args.push("-i", filePath.value)
  }
  if (workdir.value != "") {
    args.push("-w", workdir.value)
  }
  if (outputPrefix.value != "") {
    args.push("-o", outputPrefix.value)
  }
  if (noPolyA.value) {
    args.push("-long")
  }
  if (lessMem.value) {
    args.push("-lessMem")
    if (rowsLimit.value > 0) {
      args.push("-lineLimit", rowsLimit.value.toString())
    }
  }

  await RunSeqAnalysis(args).then(
    (allResult) => {
      // scrollToBottom();
      console.log(allResult);
    }
  ).catch((error) => {
    alert(error);
    console.log(error);
  });
}

onMounted(() => {
  // 订阅事件
  EventsOn("stderr-output", (data: any) => {
    // 假设你有一个用来显示输出的HTML元素
    // stderrElement.textContent += `\n${output}`;
    messages.value.push(data);
    console.log(data);
    var element = document.getElementById("log");
    if (element !== null) {
      element.scrollTop = element.scrollHeight;;
    }
    // const logElement = document.getElementById("log");
    // logElement.scrollTop = stderrElement.scrollHeight-logElement.clientHeight; // 滚动到底部
  });
});


onBeforeUnmount(() => {
  // 取消订阅事件
  EventsOff('stderr-output');
});


</script>


<template>
  <div class="container mx-auto py-4 align-baseline		" style="--wails-draggable:no-drag">
    <div class="flex w-full py-2">
      <label class="w-1/6 text-end px-4 ">{{ t("seqAnalysispage.selectFile") }}</label>
      <input class="w-5/6 px-2" type="text" @click="openFileDialog" v-model="filePath" required />
    </div>
    <div class="flex w-full py-2">
      <label class="w-1/6 text-end px-4 ">{{ t("seqAnalysispage.workdir") }}</label>
      <input class="w-5/6 px-2" type="text" v-model="workdir" required />
    </div>
    <div class="flex w-full py-2">
      <label class="w-1/6 text-end px-4 ">{{ t("seqAnalysispage.prefix") }}</label>
      <input class="w-5/6 px-2" type="text" v-model="outputPrefix" required />
    </div>
    <div class="flex py-2 justify-items-start justify-self-start justify-start">
      <div class="flex w-2/6">
        <label class="w-1/2 text-end px-4 align-bottom align-text-bottom">{{ t("seqAnalysispage.noPolyA") }}</label>
        <input class="w-1/12" type="checkbox" v-model="noPolyA" />
      </div>
      <div class="flex w-2/6">
        <label class="w-1/2 text-end px-4 align-middle	">{{ t("seqAnalysispage.lessMem") }}</label>
        <input class="w-1/12" type="checkbox" v-model="lessMem" />
      </div>
      <div class="flex w-2/6">
        <label class="w-1/2 text-end px-4 ">{{ t("seqAnalysispage.rowsLimit") }}</label>
        <input class="w-1/2 text-end px-2" type="number" v-model="rowsLimit" />
      </div>
    </div>
    <div class="flex w-full py-2 justify-end">
      <button class="button w-5/6 bg-white" @click="runSeqAnalysis">{{ t("seqAnalysispage.analysisBtn") }}</button>
    </div>
  </div>
  <div class="container mx-auto p-4 text-center">
    {{ t("seqAnalysispage.log") }}
  </div>
  <div id="log" class="log container mx-auto px-4 py-8 overflow-scroll">
    <p v-for="(line, index) in messages" :key="index">
      {{ line }}
    </p>
  </div>
</template>

<style lang="scss">
.home {
  .logo {
    display: block;
    width: 500px;
    height: 500px;
    margin: 10px auto 10px;
  }

  .link {
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    align-items: center;
    justify-content: center;
    margin: 18px auto;

    .btn {
      display: block;
      width: 150px;
      height: 50px;
      line-height: 50px;
      padding: 0 5px;
      margin: 0 30px;
      border-radius: 8px;
      text-align: center;
      font-weight: 700;
      font-size: 16px;
      white-space: nowrap;
      text-decoration: none;
      cursor: pointer;

      &.start {
        background-color: #fd0404;
        color: #ffffff;

        &:hover {
          background-color: #ec2e2e;
        }
      }

      &.star {
        background-color: #ffffff;
        color: #fd0404;

        &:hover {
          background-color: #f3f3f3;
        }
      }
    }
  }
}

.log {
  // max-width: 90vw;
  height: 50vh;
  // overflow: auto;
  text-wrap: nowrap;
  background-color: #012456;
  color: #ffffff;
}

.button {
  background-color: rgba(171, 126, 220, 0.9);

  &:hover {
    background-color: #d7a8d8;
    color: #ffffff;
  }
}

table {
  user-select: none;
  background-color: #d7a8d8;

  // th:nth-child(2),
  // th:nth-child(3),
  td:nth-child(2),
  td:nth-child(3) {
    user-select: text;
  }
}

thead tr {
  background-color: rgba(171, 126, 220, 0.9);

  &:hover {
    background-color: #d7a8d8;
    color: #ffffff;
  }
}

tbody {
  background-color: #d7a8d8;

  tr {
    background-color: #d7a8d8;

    &:hover {
      color: #ffffff;
    }

  }
}
</style>
