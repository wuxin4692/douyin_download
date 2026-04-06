import * as wails from '../bindings';

// 使用 Wails 生成的类型
const api = wails;

export const apiService = {
  // FFmpeg 状态
  getFFmpegStatus: async () => {
    return await api.GetFFmpegStatus();
  },

  // 设置
  getSettings: async () => {
    return await api.GetSettings();
  },

  saveSettings: async (settings: any) => {
    return await api.SaveSettings(settings);
  },

  // 打开目录选择对话框
  openDirectoryDialog: async (title: string) => {
    return await api.OpenDirectoryDialog(title);
  },

  // 视频解析
  parseVideo: async (url: string) => {
    const res = await api.ParseVideo(url);
    if (res.code !== 0) {
      throw new Error(res.message || '解析失败');
    }
    return { task_id: res.task_id! };
  },

  // 视频下载
  downloadVideo: async (url: string) => {
    const res = await api.DownloadVideo(url);
    if (res.code !== 0) {
      throw new Error(res.message || '下载失败');
    }
    return { task_id: res.task_id! };
  },

  // 音频提取
  extractAudio: async (videoPath: string) => {
    const res = await api.ExtractAudio(videoPath);
    if (res.code !== 0) {
      throw new Error(res.message || '音频提取失败');
    }
    return { task_id: res.task_id! };
  },

  // 语音转写
  transcribe: async (audioPath: string, apiKey?: string) => {
    const res = await api.TranscribeAudio(audioPath, apiKey);
    if (res.code !== 0) {
      throw new Error(res.message || '语音转写失败');
    }
    return { task_id: res.task_id! };
  },

  // 语义分段
  segment: async (text: string, apiKey?: string) => {
    const res = await api.SemanticSegment(text, apiKey);
    if (res.code !== 0) {
      throw new Error(res.message || '语义分段失败');
    }
    return { task_id: res.task_id! };
  },

  // LLM 分析文案
  analyzeText: async (text: string, apiKey?: string, provider?: string) => {
    const res = await api.AnalyzeText(text, apiKey, provider);
    if (res.code !== 0) {
      throw new Error(res.message || '分析失败');
    }
    return { task_id: res.task_id! };
  },

  // 一键改写
  rewriteText: async (
    text: string,
    apiKey?: string,
    provider?: string,
    style?: string,
    customInstruction?: string
  ) => {
    const res = await api.RewriteText(text, apiKey, provider, style, customInstruction);
    if (res.code !== 0) {
      throw new Error(res.message || '改写失败');
    }
    return { task_id: res.task_id! };
  },

  // 获取任务状态
  getTaskStatus: async (taskId: string) => {
    return await api.GetTaskStatus(taskId);
  },

  // 保存文案到文件
  saveFile: async (content: string, filename: string) => {
    const res = await api.SaveFile(content, filename);
    if (res.code !== 0) {
      throw new Error(res.message || '保存失败');
    }
    return { path: res.data!.path };
  },
};
