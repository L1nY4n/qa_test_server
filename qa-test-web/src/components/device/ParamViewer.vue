<template>
  <div class="param-viewer">
    <template v-if="isPrimitive(node)">
      <div class="single-value">
        <a-switch v-if="typeof node === 'boolean'" :checked="node" disabled />
        <a-input-number
          v-else-if="typeof node === 'number'"
          :value="node"
          :controls="false"
          :disabled="true"
          style="width: 180px"
        />
        <a-tag v-else-if="node === null || node === undefined || node === ''">--</a-tag>
        <a-input v-else :value="String(node)" readonly />
      </div>
    </template>

    <template v-else-if="isPrimitiveArray(node)">
      <div class="array-tags">
        <a-tag v-for="(item, idx) in primitiveArrayPreview(node)" :key="idx">
          {{ shortText(item) }}
        </a-tag>
        <a-tag v-if="primitiveArrayRemain(node) > 0" color="blue">
          +{{ primitiveArrayRemain(node) }} 项
        </a-tag>
      </div>
    </template>

    <template v-else-if="Array.isArray(node)">
      <a-collapse ghost :destroyInactivePanel="true">
        <a-collapse-panel
          v-for="(item, index) in objectArrayPreview(node)"
          :key="index"
          :header="`[${index}] ${valueHint(item)}`"
        >
          <ParamViewer :node="item" :depth="depth + 1" />
        </a-collapse-panel>
      </a-collapse>
      <a-tag v-if="objectArrayRemain(node) > 0" color="blue">
        其余 {{ objectArrayRemain(node) }} 项已折叠
      </a-tag>
    </template>

    <template v-else-if="isPlainObject(node)">
      <a-descriptions v-if="primitiveEntries.length" :column="1" size="small" bordered class="primitive-block">
        <a-descriptions-item
          v-for="[key, value] in primitiveEntries"
          :key="key"
          :label="formatKey(key)"
        >
          <template v-if="isPrimitiveArray(value)">
            <div class="array-tags">
              <a-tag v-for="(item, idx) in primitiveArrayPreview(value)" :key="idx">
                {{ shortText(item) }}
              </a-tag>
              <a-tag v-if="primitiveArrayRemain(value) > 0" color="blue">
                +{{ primitiveArrayRemain(value) }} 项
              </a-tag>
            </div>
          </template>
          <template v-else-if="typeof value === 'boolean'">
            <a-switch :checked="value" disabled />
          </template>
          <template v-else-if="typeof value === 'number'">
            <a-input-number
              :value="value"
              :controls="false"
              :disabled="true"
              style="width: 180px"
            />
          </template>
          <template v-else-if="value === null || value === undefined || value === ''">
            <a-tag>--</a-tag>
          </template>
          <template v-else>
            <a-input :value="String(value)" readonly />
          </template>
        </a-descriptions-item>
      </a-descriptions>

      <a-collapse v-if="nestedEntries.length" ghost :destroyInactivePanel="true">
        <a-collapse-panel
          v-for="[key, value] in nestedEntries"
          :key="key"
          :header="`${formatKey(key)} ${valueHint(value)}`"
        >
          <ParamViewer :node="value" :depth="depth + 1" />
        </a-collapse-panel>
      </a-collapse>
    </template>

    <template v-else>
      <a-tag>暂不支持显示该类型</a-tag>
    </template>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, PropType } from 'vue'

const ARRAY_PREVIEW_LIMIT = 20
const OBJECT_ARRAY_LIMIT = 40

function isPlainObject(value: unknown): value is Record<string, unknown> {
  return Object.prototype.toString.call(value) === '[object Object]'
}

function isPrimitive(value: unknown): value is string | number | boolean | null | undefined {
  return (
    value === null ||
    value === undefined ||
    typeof value === 'string' ||
    typeof value === 'number' ||
    typeof value === 'boolean'
  )
}

function isPrimitiveArray(value: unknown): value is Array<string | number | boolean | null | undefined> {
  return Array.isArray(value) && value.every((item) => isPrimitive(item))
}

export default defineComponent({
  name: 'ParamViewer',
  props: {
    node: {
      type: null as unknown as PropType<unknown>,
      required: true,
    },
    depth: {
      type: Number,
      default: 0,
    },
  },
  setup(props) {
    const objectEntries = computed<[string, unknown][]>(() => {
      if (!isPlainObject(props.node)) {
        return []
      }
      return Object.entries(props.node)
    })

    const primitiveEntries = computed(() =>
      objectEntries.value.filter(([, value]) => isPrimitive(value) || isPrimitiveArray(value))
    )

    const nestedEntries = computed(() =>
      objectEntries.value.filter(([, value]) => !isPrimitive(value) && !isPrimitiveArray(value))
    )

    const formatKey = (key: string): string => key.replace(/_/g, ' ')

    const shortText = (value: unknown): string => {
      if (value === null || value === undefined || value === '') {
        return '--'
      }
      const text = String(value)
      return text.length > 30 ? `${text.slice(0, 30)}...` : text
    }

    const primitiveArrayPreview = (value: unknown) => {
      if (!isPrimitiveArray(value)) {
        return []
      }
      return value.slice(0, ARRAY_PREVIEW_LIMIT)
    }

    const primitiveArrayRemain = (value: unknown): number => {
      if (!isPrimitiveArray(value)) {
        return 0
      }
      return Math.max(0, value.length - ARRAY_PREVIEW_LIMIT)
    }

    const objectArrayPreview = (value: unknown): unknown[] => {
      if (!Array.isArray(value)) {
        return []
      }
      return value.slice(0, OBJECT_ARRAY_LIMIT)
    }

    const objectArrayRemain = (value: unknown): number => {
      if (!Array.isArray(value)) {
        return 0
      }
      return Math.max(0, value.length - OBJECT_ARRAY_LIMIT)
    }

    const valueHint = (value: unknown): string => {
      if (Array.isArray(value)) {
        return `(数组 ${value.length})`
      }
      if (isPlainObject(value)) {
        return `(对象 ${Object.keys(value).length})`
      }
      if (isPrimitive(value)) {
        return `(值)`
      }
      return ''
    }

    return {
      formatKey,
      isPlainObject,
      isPrimitive,
      isPrimitiveArray,
      nestedEntries,
      objectArrayPreview,
      objectArrayRemain,
      primitiveEntries,
      primitiveArrayPreview,
      primitiveArrayRemain,
      shortText,
      valueHint,
    }
  },
})
</script>

<style scoped>
.param-viewer {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.primitive-block {
  background: #fff;
}

.array-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.single-value {
  max-width: 520px;
}
</style>
