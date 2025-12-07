import { writable, type Writable } from 'svelte/store';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import { executionState } from './executionStore';

export type ResourceType = 'META' | 'WORKER' | 'CONTAINER' | 'ORCHESTRATOR';
export type ResourceStatus = 'IDLE' | 'RUNNING' | 'THINKING' | 'PAUSED' | 'ERROR' | 'DONE';

export interface ResourceNode {
  id: string;
  name: string;
  type: ResourceType;
  status: ResourceStatus;
  detail?: string;
  children?: ResourceNode[];
  expanded?: boolean;
}

// Initial State
const initialState: ResourceNode[] = [
    {
        id: 'global-orch',
        name: 'Multiverse Orchestrator',
        type: 'ORCHESTRATOR',
        status: 'IDLE',
        expanded: true,
        children: []
    }
];

export const processResources: Writable<ResourceNode[]> = writable(initialState);

// --- Simulation Logic (Inferring from Logs) ---

// Helper to find or create a Meta Node
function getMetaNode(nodes: ResourceNode[], taskId: string): ResourceNode {
    const orch = nodes.find(n => n.id === 'global-orch');
    if (!orch) return nodes[0]; // Fallback

    const metaId = `meta-${taskId}`;
    let meta = orch.children?.find(n => n.id === metaId);

    if (!meta) {
        meta = {
            id: metaId,
            name: `Meta-Agent`, // Could be more specific if we had agent name
            type: 'META',
            status: 'IDLE',
            detail: 'Initializing...',
            expanded: true,
            children: []
        };
        orch.children = [...(orch.children || []), meta];
    }
    return meta;
}

// Helper to find or create a Worker Node
function getWorkerNode(metaNode: ResourceNode, workerId: string = 'worker-default'): ResourceNode {
    // Prefix worker ID with meta ID to ensure uniqueness globally if workerId is generic
    const uniqueWorkerId = `${metaNode.id}-${workerId}`;
    let worker = metaNode.children?.find(n => n.id === uniqueWorkerId);

    if (!worker) {
        worker = {
            id: uniqueWorkerId,
            name: `Worker: ${workerId === 'worker-default' ? 'Codex' : workerId}`,
            type: 'WORKER',
            status: 'IDLE',
            expanded: true,
            children: []
        };
        metaNode.children = [...(metaNode.children || []), worker];
    }
    return worker;
}

// Helper to find or create a Container Node
function getContainerNode(workerNode: ResourceNode, containerId: string): ResourceNode {
    let container = workerNode.children?.find(n => n.id === containerId);

    if (!container) {
        container = {
            id: containerId,
            name: `Container: ${containerId.substring(0, 8)}`,
            type: 'CONTAINER',
            status: 'IDLE',
            detail: '',
            expanded: true
        };
        workerNode.children = [...(workerNode.children || []), container];
    }
    return container;
}

export function initProcessEvents() {
    // Listen to Execution State to update Orchestrator Status
    executionState.subscribe(state => {
        processResources.update(nodes => {
            const newNodes = JSON.parse(JSON.stringify(nodes));
            const orch = newNodes.find((n: ResourceNode) => n.id === 'global-orch');
            if (orch) {
                orch.status = state === 'RUNNING' ? 'RUNNING' : (state === 'PAUSED' ? 'PAUSED' : 'IDLE');
                orch.detail = state === 'RUNNING' ? 'Orchestrating tasks...' : 'Waiting for tasks...';
            }
            return newNodes;
        });
    });

    // 1. Meta Update
    EventsOn('process:metaUpdate', (event: any) => {
        processResources.update(nodes => {
            const newNodes = JSON.parse(JSON.stringify(nodes));
            const meta = getMetaNode(newNodes, event.taskId);
            
            meta.status = event.state as ResourceStatus; // e.g. THINKING, ACTING
            meta.detail = event.detail;

            // If thinking/acting, ensure running
            if (['THINKING', 'ACTING', 'PLANNING'].includes(meta.status)) {
                const orch = newNodes.find((n: ResourceNode) => n.id === 'global-orch');
                if (orch) orch.status = 'RUNNING';
            }

            return newNodes;
        });
    });

    // 2. Worker Update
    EventsOn('process:workerUpdate', (event: any) => {
        processResources.update(nodes => {
            const newNodes = JSON.parse(JSON.stringify(nodes));
            const meta = getMetaNode(newNodes, event.taskId);
            const worker = getWorkerNode(meta, event.workerId || 'worker-default');

            worker.status = event.status as ResourceStatus;
            worker.detail = event.command || worker.detail;

            // If exit code present
            if (event.exitCode !== undefined) {
                 worker.detail += ` (Exit: ${event.exitCode})`;
            }

            return newNodes;
        });
    });

    // 3. Container Update
    EventsOn('process:containerUpdate', (event: any) => {
        processResources.update(nodes => {
            const newNodes = JSON.parse(JSON.stringify(nodes));
            const meta = getMetaNode(newNodes, event.taskId);
            const worker = getWorkerNode(meta, 'worker-default'); // Assume default worker for container
            const container = getContainerNode(worker, event.containerId);

            container.status = event.status as ResourceStatus;
            if (event.image) {
                container.detail = `Image: ${event.image}`;
            }

            return newNodes;
        });
    });

    // Validating fallback log listener (optional, but keep for now if needed)
    // Removed for cleaner Phase 2
}
