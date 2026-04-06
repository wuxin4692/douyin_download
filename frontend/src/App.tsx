import { useState, useEffect, useCallback } from 'react';
import { Settings, RefreshCw, Play, Download, Music, FileText, Wand2, Sparkles, Edit3, Check, AlertCircle } from 'lucide-react';
import { apiService } from './services/api';
import type { VideoInfo, Task, Settings as SettingsType, DownloadResult, AudioResult, TranscriptionResult, SegmentResult, LLMAnalysisResult, RewriteResult, RewriteStyle } from './types';
import { REWRITE_STYLES } from './types';
import './wailsjs/go/main/App'; // 导入 Wails 绑定

// Wails 类型声明
declare global {
  interface Window {
    go: {
      main: {
        App: {
          OpenDirectoryDialog: (title: string) => Promise<string>;
          [key: string]: any;
        };
      };
    };
  }
}

function App() {
  const [url, setUrl] = useState('');
  const [showSettings, setShowSettings] = useState(false);
  const [settings, setSettings] = useState<SettingsType>({
    silicon_flow_key: '',
    minimax_key: '',
    download_dir: '',
    ffmpeg_dir: '',
    ai_provider: 'openai',
    ai_base_url: '',
    ai_model: '',
  });
  const [ffmpegInstalled, setFfmpegInstalled] = useState<boolean | null>(null);

  const [videoInfo, setVideoInfo] = useState<VideoInfo | null>(null);
  const [currentTask, setCurrentTask] = useState<Task | null>(null);
  const [downloadResult, setDownloadResult] = useState<DownloadResult | null>(null);
  const [audioResult, setAudioResult] = useState<AudioResult | null>(null);
  const [transcription, setTranscription] = useState<TranscriptionResult | null>(null);
  const [segments, setSegments] = useState<SegmentResult | null>(null);
  const [analysis, setAnalysis] = useState<LLMAnalysisResult | null>(null);
  const [rewriteResult, setRewriteResult] = useState<RewriteResult | null>(null);
  const [selectedStyle, setSelectedStyle] = useState<RewriteStyle>('inspiring');
  const [showRewriteModal, setShowRewriteModal] = useState(false);
  const [customInstruction, setCustomInstruction] = useState('');

  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadSettings();
    checkFFmpeg();
  }, []);

  const loadSettings = async () => {
    try {
      const data = await apiService.getSettings();
      setSettings(data);
    } catch (err) {
      console.error('Failed to load settings:', err);
    }
  };

  const checkFFmpeg = async () => {
    try {
      const status = await apiService.getFFmpegStatus();
      setFfmpegInstalled(status.installed ?? false);
    } catch (err) {
      setFfmpegInstalled(false);
    }
  };

  const saveSettings = async () => {
    try {
      await apiService.saveSettings(settings);
      setShowSettings(false);
    } catch (err) {
      console.error('Failed to save settings:', err);
    }
  };

  const pollTask = useCallback(async (taskId: string, onUpdate: (task: Task) => void) => {
    const poll = async () => {
      try {
        const task = await apiService.getTaskStatus(taskId);
        onUpdate(task);
        if (task.status === 'pending' || task.status === 'processing') {
          setTimeout(poll, 500);
        }
      } catch (err) {
        console.error('Poll error:', err);
      }
    };
    poll();
  }, []);

  const handleParse = async () => {
    if (!url.trim()) return;
    setIsLoading(true);
    setError(null);

    try {
      const { task_id } = await apiService.parseVideo(url);
      pollTask(task_id, (task) => {
        setCurrentTask(task);
        if (task.status === 'completed' && task.result) {
          setVideoInfo(task.result);
        } else if (task.status === 'failed') {
          setError(task.error || '解析失败');
        }
      });
    } catch (err: any) {
      setError(err.response?.data?.message || '解析失败');
    } finally {
      setIsLoading(false);
    }
  };

  const handleDownload = async () => {
    if (!videoInfo) return;
    setIsLoading(true);
    setError(null);
    setDownloadResult(null);

    try {
      const { task_id } = await apiService.downloadVideo(url);
      pollTask(task_id, (task) => {
        setCurrentTask(task);
        if (task.status === 'completed' && task.result) {
          setDownloadResult(task.result);
        } else if (task.status === 'failed') {
          setError(task.error || '下载失败');
        }
      });
    } catch (err: any) {
      setError(err.response?.data?.message || '下载失败');
    } finally {
      setIsLoading(false);
    }
  };

  const handleExtractAudio = async () => {
    if (!downloadResult) return;
    setIsLoading(true);
    setError(null);

    try {
      const { task_id } = await apiService.extractAudio(downloadResult.video_path!);
      pollTask(task_id, (task) => {
        setCurrentTask(task);
        if (task.status === 'completed' && task.result) {
          setAudioResult(task.result);
        } else if (task.status === 'failed') {
          setError(task.error || '音频提取失败');
        }
      });
    } catch (err: any) {
      setError(err.response?.data?.message || '音频提取失败');
    } finally {
      setIsLoading(false);
    }
  };

  const handleTranscribe = async () => {
    if (!audioResult) return;
    setIsLoading(true);
    setError(null);

    try {
      const { task_id } = await apiService.transcribe(audioResult.audio_path!, settings.silicon_flow_key || undefined);
      pollTask(task_id, (task) => {
        setCurrentTask(task);
        if (task.status === 'completed' && task.result) {
          setTranscription(task.result);
        } else if (task.status === 'failed') {
          setError(task.error || '语音转写失败');
        }
      });
    } catch (err: any) {
      setError(err.response?.data?.message || '语音转写失败');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSegment = async () => {
    if (!transcription) return;
    setIsLoading(true);
    setError(null);

    try {
      const { task_id } = await apiService.segment(transcription.text!, settings.minimax_key || undefined);
      pollTask(task_id, (task) => {
        setCurrentTask(task);
        if (task.status === 'completed' && task.result) {
          setSegments(task.result);
        } else if (task.status === 'failed') {
          setError(task.error || '语义分段失败');
        }
      });
    } catch (err: any) {
      setError(err.response?.data?.message || '语义分段失败');
    } finally {
      setIsLoading(false);
    }
  };

  // LLM 分析文案
  const handleAnalyze = async () => {
    if (!transcription) return;
    setIsLoading(true);
    setError(null);
    setAnalysis(null);

    try {
      const { task_id } = await apiService.analyzeText(transcription.text!, settings.minimax_key || undefined);
      pollTask(task_id, (task) => {
        setCurrentTask(task);
        if (task.status === 'completed' && task.result) {
          setAnalysis(task.result);
        } else if (task.status === 'failed') {
          setError(task.error || '分析失败');
        }
      });
    } catch (err: any) {
      setError(err.response?.data?.message || '分析失败');
    } finally {
      setIsLoading(false);
    }
  };

  // 一键改写
  const handleRewrite = async () => {
    if (!transcription) return;
    setIsLoading(true);
    setError(null);
    setRewriteResult(null);
    setShowRewriteModal(false);

    try {
      const { task_id } = await apiService.rewriteText(
        transcription.text!,
        settings.minimax_key || undefined,
        undefined,
        selectedStyle,
        customInstruction
      );
      pollTask(task_id, (task) => {
        setCurrentTask(task);
        if (task.status === 'completed' && task.result) {
          setRewriteResult(task.result);
        } else if (task.status === 'failed') {
          setError(task.error || '改写失败');
        }
      });
    } catch (err: any) {
      setError(err.response?.data?.message || '改写失败');
    } finally {
      setIsLoading(false);
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  const handleSaveFile = async () => {
    if (!transcription) return;
    setIsLoading(true);
    setError(null);

    try {
      // 构建文件名
      const filename = videoInfo?.title 
        ? videoInfo.title.replace(/[^\w\u4e00-\u9fa5]/g, '_').substring(0, 50)
        : `douyin_${videoInfo?.video_id || Date.now()}`;
      
      await apiService.saveFile(transcription.text!, filename);
      alert('文案已保存到文件');
    } catch (err: any) {
      setError(err.response?.data?.message || '保存失败');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen p-6">
      {/* Header */}
      <header className="glass-card rounded-2xl p-6 mb-6 flex items-center justify-between animate-fade-in">
        <div className="flex items-center gap-4">
          <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-primary to-primary-hover flex items-center justify-center">
            <Play className="w-6 h-6 text-white" fill="white" />
          </div>
          <div>
            <h1 className="text-2xl font-bold gradient-text">抖音工具箱</h1>
            <p className="text-text-muted text-sm">抖音视频下载与文案提取</p>
          </div>
        </div>
        <div className="flex items-center gap-4">
          <div className={`flex items-center gap-2 px-3 py-1.5 rounded-full text-sm ${
            ffmpegInstalled ? 'bg-success/20 text-success' : 'bg-error/20 text-error'
          }`}>
            {ffmpegInstalled ? (
              <>
                <Check className="w-4 h-4" />
                <span>FFmpeg 就绪</span>
              </>
            ) : (
              <>
                <AlertCircle className="w-4 h-4" />
                <span className="cursor-pointer" onClick={() => setShowSettings(true)} title="点击设置">FFmpeg 未配置</span>
              </>
            )}
          </div>
          <button
            onClick={() => setShowSettings(true)}
            className="p-2 rounded-lg hover:bg-background-hover transition-colors"
          >
            <Settings className="w-5 h-5" />
          </button>
        </div>
      </header>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Left Panel - Input */}
        <div className="space-y-6">
          {/* URL Input */}
          <div className="glass-card rounded-2xl p-6 animate-fade-in">
            <h2 className="text-lg font-semibold mb-4">输入链接</h2>
            <div className="flex gap-3">
              <input
                type="text"
                value={url}
                onChange={(e) => setUrl(e.target.value)}
                placeholder="粘贴抖音分享链接或视频ID..."
                className="flex-1 bg-background border border-white/10 rounded-xl px-4 py-3 text-text placeholder:text-text-muted focus:outline-none focus:border-primary transition-colors"
              />
              <button
                onClick={handleParse}
                disabled={isLoading || !url.trim()}
                className="btn-gradient px-6 py-3 rounded-xl font-medium disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
              >
                <Play className="w-4 h-4" />
                解析
              </button>
            </div>
          </div>

          {/* Video Info */}
          {videoInfo && (
            <div className="glass-card rounded-2xl p-6 animate-slide-up">
              <h2 className="text-lg font-semibold mb-4">视频信息</h2>
              <div className="flex gap-4">
                {videoInfo.cover_url && (
                  <img
                    src={videoInfo.cover_url}
                    alt={videoInfo.title}
                    className="w-32 h-32 object-cover rounded-xl"
                  />
                )}
                <div className="flex-1">
                  <h3 className="font-medium text-lg mb-2 line-clamp-2">{videoInfo.title || '无标题'}</h3>
                  <p className="text-text-muted text-sm mb-1">作者: {videoInfo.author || '未知'}</p>
                  <p className="text-text-muted text-sm">ID: {videoInfo.video_id}</p>
                </div>
              </div>
              <div className="mt-4 flex gap-3">
                <button
                  onClick={handleDownload}
                  disabled={isLoading || !videoInfo}
                  className="btn-gradient px-4 py-2 rounded-lg font-medium flex items-center gap-2 disabled:opacity-50"
                >
                  <Download className="w-4 h-4" />
                  {downloadResult ? '重新下载' : '下载视频'}
                </button>
                {downloadResult && (
                  <>
                    <button
                      onClick={handleExtractAudio}
                      disabled={isLoading || !downloadResult}
                      className="bg-background-hover px-4 py-2 rounded-lg font-medium flex items-center gap-2 hover:bg-background transition-colors disabled:opacity-50"
                    >
                      <Music className="w-4 h-4" />
                      提取音频
                    </button>
                    <button
                      onClick={handleTranscribe}
                      disabled={isLoading || !audioResult}
                      className="bg-background-hover px-4 py-2 rounded-lg font-medium flex items-center gap-2 hover:bg-background transition-colors disabled:opacity-50"
                    >
                      <FileText className="w-4 h-4" />
                      语音转写
                    </button>
                    <button
                      onClick={handleSegment}
                      disabled={isLoading || !transcription}
                      className="bg-background-hover px-4 py-2 rounded-lg font-medium flex items-center gap-2 hover:bg-background transition-colors disabled:opacity-50"
                    >
                      <Wand2 className="w-4 h-4" />
                      语义分段
                    </button>
                    {transcription && (
                      <>
                        <button
                          onClick={handleAnalyze}
                          disabled={isLoading}
                          className="bg-purple-500/20 text-purple-400 px-4 py-2 rounded-lg font-medium flex items-center gap-2 hover:bg-purple-500/30 transition-colors disabled:opacity-50"
                        >
                          <Sparkles className="w-4 h-4" />
                          AI分析
                        </button>
                        <button
                          onClick={() => setShowRewriteModal(true)}
                          disabled={isLoading}
                          className="bg-cyan-500/20 text-cyan-400 px-4 py-2 rounded-lg font-medium flex items-center gap-2 hover:bg-cyan-500/30 transition-colors disabled:opacity-50"
                        >
                          <Edit3 className="w-4 h-4" />
                          一键改写
                        </button>
                      </>
                    )}
                  </>
                )}
              </div>
            </div>
          )}

          {/* Progress */}
          {currentTask && currentTask.status !== 'completed' && (
            <div className="glass-card rounded-2xl p-6 animate-slide-up">
              <h2 className="text-lg font-semibold mb-4">处理进度</h2>
              <div className="mb-2 flex justify-between text-sm">
                <span>{currentTask.message}</span>
                <span>{currentTask.progress}%</span>
              </div>
              <div className="h-2 bg-background rounded-full overflow-hidden">
                <div
                  className="h-full bg-gradient-to-r from-primary to-primary-hover transition-all duration-300"
                  style={{ width: `${currentTask.progress}%` }}
                />
              </div>
            </div>
          )}

          {/* Error */}
          {error && (
            <div className="glass-card rounded-2xl p-6 border border-error/50 animate-slide-up">
              <p className="text-error">{error}</p>
            </div>
          )}
        </div>

        {/* Right Panel - Results */}
        <div className="space-y-6">
          {/* Transcription */}
          {transcription && (
            <div className="glass-card rounded-2xl p-6 animate-slide-up">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-lg font-semibold">语音转写</h2>
                <div className="flex gap-3">
                  <button
                    onClick={handleSaveFile}
                    disabled={isLoading}
                    className="text-sm text-primary hover:text-primary-hover disabled:opacity-50"
                  >
                    保存到文件
                  </button>
                <button
                  onClick={() => copyToClipboard(transcription.text || '')}
                  className="text-sm text-primary hover:text-primary-hover"
                >
                  复制全文
                </button>
                </div>
              </div>
              <div className="bg-background rounded-xl p-4 max-h-64 overflow-y-auto">
                <p className="text-text whitespace-pre-wrap">{transcription.text}</p>
              </div>
            </div>
          )}

          {/* Segments */}
          {segments && (
            <div className="glass-card rounded-2xl p-6 animate-slide-up">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-lg font-semibold">语义分段</h2>
                <button
                  onClick={() => copyToClipboard(
                    (segments.segments || []).map(s => `## ${s.title}\n\n${s.content}`).join('\n\n')
                  )}
                  className="text-sm text-primary hover:text-primary-hover"
                >
                  复制全部
                </button>
              </div>
              <div className="space-y-4 max-h-96 overflow-y-auto">
                {segments.segments?.map((segment, index) => (
                  <div key={index} className="bg-background rounded-xl p-4">
                    <h3 className="font-medium text-primary mb-2">## {segment.title}</h3>
                    <p className="text-text-muted whitespace-pre-wrap">{segment.content}</p>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* LLM Analysis */}
          {analysis && (
            <div className="glass-card rounded-2xl p-6 animate-slide-up border border-purple-500/30">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-lg font-semibold flex items-center gap-2">
                  <Sparkles className="w-5 h-5 text-purple-400" />
                  AI 分析报告
                </h2>
                <button
                  onClick={() => copyToClipboard(analysis.analysis || '')}
                  className="text-sm text-primary hover:text-primary-hover"
                >
                  复制
                </button>
              </div>
              <div className="bg-background rounded-xl p-4 max-h-80 overflow-y-auto">
                <p className="text-text whitespace-pre-wrap leading-relaxed">{analysis.analysis || ''}</p>
              </div>
            </div>
          )}

          {/* Rewrite */}
          {rewriteResult && (
            <div className="glass-card rounded-2xl p-6 animate-slide-up border border-cyan-500/30">
              <div className="flex items-center justify-between mb-4">
                  <h2 className="text-lg font-semibold flex items-center gap-2">
                    <Edit3 className="w-5 h-5 text-cyan-400" />
                    一键改写 ({REWRITE_STYLES.find(s => s.value === rewriteResult.style)?.label || rewriteResult.style || ''})
                  </h2>
                <button
                  onClick={() => copyToClipboard(rewriteResult.rewritten || '')}
                  className="text-sm text-primary hover:text-primary-hover"
                >
                  复制改写后文案
                </button>
              </div>

              {/* Style Selector */}
              <div className="mb-4">
                <label className="block text-sm text-text-muted mb-2">改写风格</label>
                <div className="flex flex-wrap gap-2">
                  {REWRITE_STYLES.map((style) => (
                    <button
                      key={style.value}
                      onClick={() => {
                        setSelectedStyle(style.value);
                        setRewriteResult(null); // 重新改写
                      }}
                      disabled={isLoading}
                      className={`px-3 py-1.5 rounded-lg text-sm flex items-center gap-1 transition-colors disabled:opacity-50 ${
                        selectedStyle === style.value
                          ? 'bg-cyan-500/30 text-cyan-400'
                          : 'bg-background-hover text-text-muted hover:text-text'
                      }`}
                    >
                      <span>{style.emoji}</span>
                      <span>{style.label}</span>
                    </button>
                  ))}
                </div>
              </div>

              {/* Original Text */}
              <div className="mb-4">
                <h3 className="text-sm text-text-muted mb-2">原文</h3>
                <div className="bg-background/50 rounded-xl p-3 text-sm text-text-muted max-h-32 overflow-y-auto">
                  <p className="whitespace-pre-wrap">{rewriteResult.original || ''}</p>
                </div>
              </div>

              {/* Rewritten Text */}
              <div>
                <h3 className="text-sm text-text-muted mb-2">改写后</h3>
                <div className="bg-background rounded-xl p-4 max-h-48 overflow-y-auto">
                  <p className="text-text whitespace-pre-wrap leading-relaxed">{rewriteResult.rewritten || ''}</p>
                </div>
              </div>
            </div>
          )}

          {/* Empty State */}
          {!transcription && !segments && (
            <div className="glass-card rounded-2xl p-12 flex flex-col items-center justify-center text-center animate-fade-in">
              <div className="w-20 h-20 rounded-full bg-background-hover flex items-center justify-center mb-6">
                <FileText className="w-10 h-10 text-text-muted" />
              </div>
              <h3 className="text-xl font-semibold mb-2">等待处理</h3>
              <p className="text-text-muted max-w-xs">
                输入抖音链接后，按照步骤：下载视频 → 提取音频 → 语音转写 → 语义分段 → AI分析/一键改写
              </p>
            </div>
          )}
        </div>
      </div>

      {/* Settings Modal */}
      {showSettings && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="glass-card rounded-2xl p-6 w-full max-w-md animate-slide-up">
            <h2 className="text-xl font-semibold mb-6">设置</h2>
            <div className="space-y-4">
              <div>
                <label className="block text-sm text-text-muted mb-2">FFmpeg 目录</label>
                <div className="flex gap-2">
                  <input
                    type="text"
                    value={settings.ffmpeg_dir || ''}
                    onChange={(e) => setSettings({ ...settings, ffmpeg_dir: e.target.value })}
                    placeholder="选择 ffmpeg 所在目录"
                    className="flex-1 bg-background border border-white/10 rounded-xl px-4 py-3 text-text placeholder:text-text-muted focus:outline-none focus:border-primary"
                  />
                  <button
                    onClick={async () => {
                      try {
                        const dir = await window.go.main.App.OpenDirectoryDialog('选择 FFmpeg 目录');
                        if (dir) {
                          setSettings({ ...settings, ffmpeg_dir: dir });
                        }
                      } catch (e) {
                        console.error('Failed to open folder dialog:', e);
                      }
                    }}
                    className="px-4 py-3 rounded-xl bg-background-hover hover:bg-background-active transition-colors text-sm"
                  >
                    浏览
                  </button>
                </div>
              </div>
              <div>
                <label className="block text-sm text-text-muted mb-2">Silicon Flow API Key</label>
                <input
                  type="password"
                  value={settings.silicon_flow_key}
                  onChange={(e) => setSettings({ ...settings, silicon_flow_key: e.target.value })}
                  placeholder="用于语音转文字"
                  className="w-full bg-background border border-white/10 rounded-xl px-4 py-3 text-text placeholder:text-text-muted focus:outline-none focus:border-primary"
                />
              </div>
              <div>
                <label className="block text-sm text-text-muted mb-2">AI 模型</label>
                <select
                  value={settings.ai_provider || 'openai'}
                  onChange={(e) => setSettings({ ...settings, ai_provider: e.target.value as any })}
                  className="w-full bg-background border border-white/10 rounded-xl px-4 py-3 text-text focus:outline-none focus:border-primary"
                >
                  <option value="openai">OpenAI</option>
                  <option value="minimax">MiniMax</option>
                  <option value="compatible">OpenAI 兼容接口</option>
                </select>
              </div>
              {(settings.ai_provider === 'minimax' || settings.ai_provider === 'openai') && (
                <div>
                  <label className="block text-sm text-text-muted mb-2">具体模型</label>
                  <select
                    value={settings.ai_model || ''}
                    onChange={(e) => setSettings({ ...settings, ai_model: e.target.value })}
                    className="w-full bg-background border border-white/10 rounded-xl px-4 py-3 text-text focus:outline-none focus:border-primary"
                  >
                    {settings.ai_provider === 'openai' && (
                      <>
                        <option value="">默认 (gpt-3.5-turbo)</option>
                        <option value="gpt-3.5-turbo">GPT-3.5 Turbo</option>
                        <option value="gpt-4">GPT-4</option>
                        <option value="gpt-4-turbo-preview">GPT-4 Turbo</option>
                      </>
                    )}
                    {settings.ai_provider === 'minimax' && (
                      <>
                        <option value="">默认 (MiniMax-M2)</option>
                        <option value="MiniMax-M2">MiniMax-M2</option>
                        <option value="MiniMax-M2.5">MiniMax-M2.5</option>
                      </>
                    )}
                  </select>
                </div>
              )}
              <div>
                <label className="block text-sm text-text-muted mb-2">API Key</label>
                <input
                  type="password"
                  value={settings.minimax_key}
                  onChange={(e) => setSettings({ ...settings, minimax_key: e.target.value })}
                  placeholder="输入 API Key"
                  className="w-full bg-background border border-white/10 rounded-xl px-4 py-3 text-text placeholder:text-text-muted focus:outline-none focus:border-primary"
                />
              </div>
              {(settings.ai_provider === 'compatible') && (
                <>
                  <div>
                    <label className="block text-sm text-text-muted mb-2">API 接口地址</label>
                    <input
                      type="text"
                      value={settings.ai_base_url || ''}
                      onChange={(e) => setSettings({ ...settings, ai_base_url: e.target.value })}
                      placeholder="https://api.openai.com/v1"
                      className="w-full bg-background border border-white/10 rounded-xl px-4 py-3 text-text placeholder:text-text-muted focus:outline-none focus:border-primary"
                    />
                  </div>
                  <div>
                    <label className="block text-sm text-text-muted mb-2">模型名称</label>
                    <input
                      type="text"
                      value={settings.ai_model || ''}
                      onChange={(e) => setSettings({ ...settings, ai_model: e.target.value })}
                      placeholder="如: gpt-3.5-turbo, gpt-4, claude-3-sonnet"
                      className="w-full bg-background border border-white/10 rounded-xl px-4 py-3 text-text placeholder:text-text-muted focus:outline-none focus:border-primary"
                    />
                  </div>
                </>
              )}
            </div>
            <div className="flex gap-3 mt-6">
              <button
                onClick={() => setShowSettings(false)}
                className="flex-1 px-4 py-3 rounded-xl border border-white/10 hover:bg-background-hover transition-colors"
              >
                取消
              </button>
              <button
                onClick={saveSettings}
                className="flex-1 btn-gradient px-4 py-3 rounded-xl font-medium"
              >
                保存
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Rewrite Modal */}
      {showRewriteModal && (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="glass-card rounded-2xl p-6 w-full max-w-lg animate-slide-up">
            <h2 className="text-xl font-semibold mb-4 flex items-center gap-2">
              <Edit3 className="w-5 h-5 text-cyan-400" />
              一键改写
            </h2>

            {/* 风格选择 */}
            <div className="mb-4">
              <label className="block text-sm text-text-muted mb-2">选择风格</label>
              <div className="flex flex-wrap gap-2">
                {REWRITE_STYLES.map((style) => (
                  <button
                    key={style.value}
                    onClick={() => setSelectedStyle(style.value)}
                    className={`px-3 py-1.5 rounded-lg text-sm flex items-center gap-1 transition-colors ${
                      selectedStyle === style.value
                        ? 'bg-cyan-500/30 text-cyan-400'
                        : 'bg-background-hover text-text-muted hover:text-text'
                    }`}
                  >
                    <span>{style.emoji}</span>
                    <span>{style.label}</span>
                  </button>
                ))}
              </div>
            </div>

            {/* 自定义指令 */}
            <div className="mb-4">
              <label className="block text-sm text-text-muted mb-2">
                你的改写要求（可选）
              </label>
              <textarea
                value={customInstruction}
                onChange={(e) => setCustomInstruction(e.target.value)}
                placeholder="例如：&#10;• 加入一些幽默元素&#10;• 突出产品的性价比&#10;• 添加行动号召语&#10;• 控制在200字以内&#10;• 使用更接地气的表达..."
                className="w-full h-32 bg-background border border-white/10 rounded-xl px-4 py-3 text-text placeholder:text-text-muted/50 focus:outline-none focus:border-cyan-500/50 resize-none"
              />
              <p className="text-xs text-text-muted mt-1">
                不填则按所选风格自动改写
              </p>
            </div>

            {/* 快捷模板 */}
            <div className="mb-4">
              <label className="block text-sm text-text-muted mb-2">快捷模板</label>
              <div className="flex flex-wrap gap-2">
                {[
                  { label: '加幽默', text: '加入幽默元素，让内容更有趣' },
                  { label: '加号召', text: '结尾添加行动号召，引导用户点赞关注' },
                  { label: '简短', text: '精简内容，突出核心要点' },
                  { label: '卖货风', text: '突出产品卖点，强调优惠和价值' },
                ].map((template) => (
                  <button
                    key={template.label}
                    onClick={() => setCustomInstruction(prev => 
                      prev ? prev + '\n• ' + template.text : '• ' + template.text
                    )}
                    className="px-2 py-1 text-xs bg-background-hover rounded-lg text-text-muted hover:text-text transition-colors"
                  >
                    {template.label}
                  </button>
                ))}
              </div>
            </div>

            <div className="flex gap-3">
              <button
                onClick={() => {
                  setShowRewriteModal(false);
                  setCustomInstruction('');
                }}
                className="flex-1 px-4 py-3 rounded-xl border border-white/10 hover:bg-background-hover transition-colors"
              >
                取消
              </button>
              <button
                onClick={handleRewrite}
                disabled={isLoading}
                className="flex-1 btn-gradient px-4 py-3 rounded-xl font-medium disabled:opacity-50 flex items-center justify-center gap-2"
              >
                {isLoading ? (
                  <>
                    <RefreshCw className="w-4 h-4 animate-spin" />
                    改写中...
                  </>
                ) : (
                  <>
                    <Sparkles className="w-4 h-4" />
                    开始改写
                  </>
                )}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default App;
