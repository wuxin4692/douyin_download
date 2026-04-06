// 使用 any 类型避免与 Wails 生成类型冲突
export interface VideoInfo {
  video_id?: string;
  title?: string;
  cover_url?: string;
  download_url?: string;
  author?: string;
  duration?: number;
}

export interface Task {
  id?: string;
  status?: string;
  progress?: number;
  message?: string;
  error?: string;
  result?: any;
  created_at?: string;
  completed_at?: string;
}

export interface Settings {
  silicon_flow_key?: string;
  minimax_key?: string;
  download_dir?: string;
  ffmpeg_dir?: string;
  ai_provider?: 'openai' | 'minimax' | 'compatible';
  ai_base_url?: string;
  ai_model?: string;
}

export interface DownloadResult {
  video_path?: string;
  title?: string;
  size?: number;
}

export interface AudioResult {
  audio_path?: string;
  duration?: number;
}

export interface TranscriptionResult {
  text?: string;
  language?: string;
  duration?: number;
}

export interface Segment {
  title?: string;
  content?: string;
}

export interface SegmentResult {
  segments?: Segment[];
  full_text?: string;
}

export interface FFmpegStatus {
  installed?: boolean;
  path?: string;
}

export interface ApiResponse<T = any> {
  code: number;
  message?: string;
  data?: T;
  task_id?: string;
}

// LLM 分析结果
export interface LLMAnalysisResult {
  analysis?: string;
}

// LLM 改写结果
export interface RewriteResult {
  original?: string;
  rewritten?: string;
  style?: string;
}

// 改写风格选项
export type RewriteStyle = 'inspiring' | 'humorous' | 'professional' | 'casual' | 'emotional';

export const REWRITE_STYLES: { value: RewriteStyle; label: string; emoji: string }[] = [
  { value: 'inspiring', label: '激励人心', emoji: '💪' },
  { value: 'humorous', label: '幽默风趣', emoji: '😄' },
  { value: 'professional', label: '专业严谨', emoji: '💼' },
  { value: 'casual', label: '日常随意', emoji: '☕' },
  { value: 'emotional', label: '情感丰富', emoji: '❤️' },
];
