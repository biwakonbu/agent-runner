<script lang="ts">
  import WBSGraphNode from "./WBSGraphNode.svelte";
  import type { WBSNode } from "../../stores/wbsStore";
  import type { TaskStatus } from "../../schemas";

  // WBSNode properties exposed as controls
  export let label: string = "Sample Node";
  export let type: "phase" | "task" = "task";
  export let phaseName: string = "実装";
  export let status: TaskStatus = "PENDING";
  export let hasChildren: boolean = false;
  export let expanded: boolean = false;

  // Construct node object from flat props
  $: node = {
    id: "sample-id",
    type,
    label,
    phaseName,
    level: 1,
    children: hasChildren ? ["child-1"] : [],
    task:
      type === "task"
        ? {
            id: "task-1",
            title: label,
            status,
            phaseName,
            poolId: "default",
            createdAt: "",
            updatedAt: "",
            dependencies: [],
          }
        : undefined,
    progress: { total: 0, completed: 0, percentage: 0 },
  } as unknown as WBSNode; // Cast as WBSNode
</script>

<div style="position: relative; width: 300px; height: 100px;">
  <WBSGraphNode {node} x={50} y={20} />
</div>
