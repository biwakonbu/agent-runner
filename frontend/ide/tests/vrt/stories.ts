/**
 * Storybook ストーリー取得ユーティリティ
 *
 * Storybook の index.json から全ストーリーを取得する
 */
import { request } from '@playwright/test';

export interface StoryEntry {
  id: string;
  title: string;
  name: string;
  importPath: string;
  tags?: string[];
}

/**
 * Storybook から全ストーリーを取得
 */
export async function getStories(baseUrl = 'http://localhost:6006'): Promise<StoryEntry[]> {
  const context = await request.newContext();
  const response = await context.get(`${baseUrl}/index.json`);

  if (!response.ok()) {
    throw new Error(`Failed to fetch stories: ${response.status()}`);
  }

  const data = await response.json();

  // entries からストーリータイプのみをフィルタ
  const stories = Object.values(data.entries as Record<string, StoryEntry & { type: string }>)
    .filter((entry) => entry.type === 'story')
    .map((entry) => ({
      id: entry.id,
      title: entry.title,
      name: entry.name,
      importPath: entry.importPath,
      tags: entry.tags,
    }));

  return stories;
}

/**
 * ストーリーの iframe URL を生成
 */
export function getStoryUrl(storyId: string, baseUrl = 'http://localhost:6006'): string {
  return `${baseUrl}/iframe.html?id=${storyId}&viewMode=story`;
}
