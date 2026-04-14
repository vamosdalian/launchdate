import { useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Rocket, Calendar, Newspaper, MapPin, Building2, X, Layers, Target } from 'lucide-react';
import { statsService, taskService } from '@/services';
import type { DashboardStats } from '@/services/statsService';
import { Button } from '@/components/ui/button';
import type { TaskAction, TaskInfo, TaskProgress, TaskType } from '@/services/taskService';

const TASK_TYPE_OPTIONS: Array<{ type: TaskType; label: string }> = [
  { type: 'launch', label: 'Launch' },
  { type: 'agency', label: 'Agency' },
  { type: 'launcher', label: 'Launcher' },
  { type: 'launcher_family', label: 'Launcher Family' },
  { type: 'pad', label: 'Pad' },
  { type: 'location', label: 'Location' },
  { type: 'update', label: 'Update' },
];

function formatTaskDate(value?: string) {
  if (!value) {
    return '—';
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return date.toLocaleString();
}

function getProgressValue<T extends keyof TaskProgress>(
  progress: TaskProgress | undefined,
  key: T,
): TaskProgress[T] | undefined {
  return progress?.[key];
}

function getTaskStatusClasses(status: string) {
  switch (status) {
    case 'running':
      return 'bg-green-100 text-green-800';
    case 'paused':
      return 'bg-amber-100 text-amber-800';
    case 'completed':
      return 'bg-blue-100 text-blue-800';
    case 'canceled':
      return 'bg-slate-200 text-slate-700';
    case 'failed':
      return 'bg-red-100 text-red-800';
    default:
      return 'bg-muted text-muted-foreground';
  }
}

export default function Dashboard() {
  const navigate = useNavigate();
  const [stats, setStats] = useState<DashboardStats>({
    totalRockets: 0,
    upcomingLaunches: 0,
    newsArticles: 0,
    launchBases: 0,
    companies: 0,
    ll2: {
      launches: 0,
      agencies: 0,
      launchers: 0,
      launcher_families: 0,
      locations: 0,
      pads: 0,
    },
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [task, setTask] = useState<TaskInfo | null>(null);
  const [taskHistory, setTaskHistory] = useState<TaskInfo[]>([]);
  const [loadingTask, setLoadingTask] = useState(false);
  const [taskPendingAction, setTaskPendingAction] = useState<string | null>(null);
  const [showWelcome, setShowWelcome] = useState(
    () => localStorage.getItem('launchdate-admin-welcome-closed') !== 'true'
  );

  const fetchTask = useCallback(async () => {
    try {
      setLoadingTask(true);
      setError(null);
      const [currentTask, recentTasks] = await Promise.all([
        taskService.getCurrentTask(),
        taskService.getTaskHistory(10),
      ]);
      setTask(currentTask);
      setTaskHistory(recentTasks);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch task info');
      setTask(null);
      setTaskHistory([]);
    } finally {
      setLoadingTask(false);
    }
  }, []);

  const fetchStats = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await statsService.getStats();
      setStats(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch statistics');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchStats();
    const timer = window.setInterval(fetchStats, 30000);
    return () => window.clearInterval(timer);
  }, [fetchStats]);

  useEffect(() => {
    fetchTask();
    const timer = window.setInterval(fetchTask, 2000);
    return () => window.clearInterval(timer);
  }, [fetchTask]);

  const handleTaskAction = async (action: TaskAction) => {
    try {
      setTaskPendingAction(action);
      setError(null);
      await taskService.actionTask(action);
      await Promise.all([fetchTask(), fetchStats()]);
    } catch (err) {
      setError(err instanceof Error ? err.message : `Failed to ${action} task`);
    } finally {
      setTaskPendingAction(null);
    }
  };

  const handleStartTask = async (type: TaskType) => {
    try {
      setTaskPendingAction(`start-${type}`);
      setError(null);
      await taskService.startTask(type);
      await Promise.all([fetchTask(), fetchStats()]);
    } catch (err) {
      setError(err instanceof Error ? err.message : `Failed to start ${type} task`);
    } finally {
      setTaskPendingAction(null);
    }
  };

  const statItems = [
    {
      label: 'Total Rockets',
      value: stats.totalRockets,
      icon: Rocket,
      color: 'text-blue-600',
      path: '/rockets',
    },
    {
      label: 'Total Launches',
      value: stats.upcomingLaunches,
      icon: Calendar,
      color: 'text-green-600',
      path: '/launches',
    },
    {
      label: 'Launch Bases',
      value: stats.launchBases,
      icon: MapPin,
      color: 'text-orange-600',
      path: '/launch-bases',
    },
    {
      label: 'Agencies',
      value: stats.companies,
      icon: Building2,
      color: 'text-red-600',
      path: '/companies',
    },
    {
      label: 'News Articles',
      value: stats.newsArticles,
      icon: Newspaper,
      color: 'text-purple-600',
      path: '/news',
    },
  ];

  const ll2StatItems = [
    {
      label: 'LL2 Launchers',
      value: stats.ll2.launchers,
      icon: Rocket, // Reuse Rocket for now, or find another
      color: 'text-orange-600',
      path: '/ll2/launchers',
    },
    {
      label: 'LL2 Launches',
      value: stats.ll2.launches,
      icon: Rocket,
      color: 'text-blue-600',
      path: '/ll2/launches',
    },
    {
      label: 'LL2 Locations',
      value: stats.ll2.locations,
      icon: MapPin,
      color: 'text-green-600',
      path: '/ll2/locations',
    },
    {
      label: 'LL2 Agencies',
      value: stats.ll2.agencies,
      icon: Building2,
      color: 'text-red-600',
      path: '/ll2/agencies',
    },
    {
      label: 'LL2 Families',
      value: stats.ll2.launcher_families,
      icon: Layers,
      color: 'text-purple-600',
      path: '/ll2/families', // Verify path? Assume standard pattern for now
    },
    {
      label: 'LL2 Pads',
      value: stats.ll2.pads,
      icon: Target,
      color: 'text-pink-600',
      path: '/ll2/pads',
    },
  ];

  const taskProgress = task?.progress;
  const currentCount = Number(getProgressValue(taskProgress, 'current_count') ?? 0);
  const totalCount = Number(getProgressValue(taskProgress, 'total_count') ?? 0);
  const hasCountProgress = currentCount > 0 || totalCount > 0;

  const taskCountDetails = hasCountProgress
    ? [
        {
          label: 'Fetched',
          value: totalCount > 0 ? `${currentCount} / ${totalCount}` : String(currentCount),
        },
        ...(totalCount > 0
          ? [
              {
                label: 'Remaining',
                value: String(Math.max(totalCount - currentCount, 0)),
              },
            ]
          : []),
      ]
    : [];

  const updateTaskDetails = task?.type === 'update'
    ? [
        {
          label: 'Next Run At',
          value: formatTaskDate(getProgressValue(taskProgress, 'next_run_at')),
        },
        {
          label: 'Watermark',
          value: formatTaskDate(getProgressValue(taskProgress, 'watermark_last_updated')),
        },
        {
          label: 'Last Success',
          value: formatTaskDate(getProgressValue(taskProgress, 'last_success_at')),
        },
      ]
    : [];

  const recentTasks = task
    ? taskHistory.filter(
        (item) => !(item.type === task.type && item.started_at === task.started_at && item.updated_at === task.updated_at)
      )
    : taskHistory;

  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-8">Dashboard</h1>

      {error && (
        <div className="mb-4 p-4 bg-red-100 text-red-800 rounded-md">
          {error}
        </div>
      )}

      {loading ? (
        <div className="text-center p-8">Loading dashboard statistics...</div>
      ) : (
        <>
          {showWelcome && (
            <Card className="mb-8">
              <CardHeader className="flex flex-row items-center justify-between">
                <CardTitle>Welcome to LaunchDate Admin</CardTitle>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => {
                    setShowWelcome(false);
                    localStorage.setItem('launchdate-admin-welcome-closed', 'true');
                  }}
                >
                  <X className="h-4 w-4" />
                </Button>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">
                  Manage your rocket launch data, news, companies, and launch bases from this admin panel.
                  Use the sidebar navigation to access different sections.
                </p>
              </CardContent>
            </Card>
          )}

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-6 mb-8">
            {statItems.map((stat) => {
              const Icon = stat.icon;
              return (
                <Card
                  key={stat.label}
                  className="cursor-pointer hover:shadow-lg transition-shadow"
                  onClick={() => navigate(stat.path)}
                >
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">{stat.label}</CardTitle>
                    <Icon className={`h-4 w-4 ${stat.color}`} />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{stat.value}</div>
                  </CardContent>
                </Card>
              );
            })}
          </div>

          <h2 className="text-2xl font-bold mb-4">LL2 Statistics</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6 gap-6 mb-8">
            {ll2StatItems.map((stat) => {
              const Icon = stat.icon;
              return (
                <Card
                  key={stat.label}
                >
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">{stat.label}</CardTitle>
                    <Icon className={`h-4 w-4 ${stat.color}`} />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{stat.value}</div>
                  </CardContent>
                </Card>
              );
            })}
          </div>
          <div className="grid grid-cols-1 gap-6 xl:grid-cols-[minmax(0,1.15fr)_minmax(0,0.85fr)]">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between">
                <CardTitle>Current Task</CardTitle>
                <div className="flex items-center gap-3">
                  {loadingTask && task && (
                    <span className="text-xs text-muted-foreground">更新中...</span>
                  )}
                  <Button variant="outline" size="sm" onClick={fetchTask} disabled={loadingTask}>
                    刷新
                  </Button>
                </div>
              </CardHeader>
              <CardContent className="space-y-4">
                {loadingTask && !task ? (
                  <div className="text-sm text-muted-foreground">Loading task...</div>
                ) : task ? (
                  <div className="space-y-3">
                    <div className="grid gap-3 md:grid-cols-3">
                      <div>
                        <p className="text-xs text-muted-foreground">Type</p>
                        <p className="font-medium">{task.type}</p>
                      </div>
                      <div>
                        <p className="text-xs text-muted-foreground">Status</p>
                        <span
                          className={`inline-flex rounded-full px-2.5 py-1 text-xs font-medium ${getTaskStatusClasses(task.status)}`}
                        >
                          {task.status}
                        </span>
                      </div>
                      <div>
                        <p className="text-xs text-muted-foreground">Started At</p>
                        <p className="font-medium">{formatTaskDate(task.started_at)}</p>
                      </div>
                      <div>
                        <p className="text-xs text-muted-foreground">Updated At</p>
                        <p className="font-medium">{formatTaskDate(task.updated_at)}</p>
                      </div>
                    </div>
                    {task.finished_at && (
                      <div>
                        <p className="text-xs text-muted-foreground">Finished At</p>
                        <p className="font-medium">{formatTaskDate(task.finished_at)}</p>
                      </div>
                    )}
                    {taskCountDetails.length > 0 && (
                      <div>
                        <p className="text-xs text-muted-foreground mb-2">Progress</p>
                        <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-4">
                          {taskCountDetails.map((detail) => (
                            <div key={detail.label} className="rounded-md border bg-muted/30 p-3">
                              <p className="text-xs text-muted-foreground">{detail.label}</p>
                              <p className="font-medium break-all">{detail.value}</p>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}
                    {updateTaskDetails.length > 0 && (
                      <div>
                        <p className="text-xs text-muted-foreground mb-2">Update Sync</p>
                        <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-4">
                          {updateTaskDetails.map((detail) => (
                            <div key={detail.label} className="rounded-md border bg-muted/30 p-3">
                              <p className="text-xs text-muted-foreground">{detail.label}</p>
                              <p className="font-medium break-all">{detail.value}</p>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}
                    {task.last_error && (
                      <div className="rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">
                        {task.last_error}
                      </div>
                    )}
                    <div className="flex flex-wrap gap-2">
                      <Button
                        variant="outline"
                        size="sm"
                        disabled={taskPendingAction !== null}
                        onClick={() => handleTaskAction('pause')}
                      >
                        Pause
                      </Button>
                      <Button
                        variant="outline"
                        size="sm"
                        disabled={taskPendingAction !== null}
                        onClick={() => handleTaskAction('resume')}
                      >
                        Resume
                      </Button>
                      <Button
                        variant="destructive"
                        size="sm"
                        disabled={taskPendingAction !== null}
                        onClick={() => handleTaskAction('cancel')}
                      >
                        Cancel
                      </Button>
                    </div>
                  </div>
                ) : (
                  <div className="text-sm text-muted-foreground">No active task.</div>
                )}

                <div>
                  <p className="mb-2 text-sm font-medium">Start New Task</p>
                  <div className="flex flex-wrap gap-2">
                    {TASK_TYPE_OPTIONS.map((item) => (
                      <Button
                        key={item.type}
                        variant="outline"
                        size="sm"
                        disabled={taskPendingAction !== null}
                        onClick={() => handleStartTask(item.type)}
                      >
                        {item.label}
                      </Button>
                    ))}
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Recent Tasks</CardTitle>
              </CardHeader>
              <CardContent>
                {recentTasks.length > 0 ? (
                  <div className="space-y-2">
                    {recentTasks.map((item, index) => {
                      const historyCurrentCount = Number(item.progress?.current_count ?? 0);
                      const historyTotalCount = Number(item.progress?.total_count ?? 0);
                      const historyProgress = historyTotalCount > 0
                        ? `${historyCurrentCount} / ${historyTotalCount}`
                        : historyCurrentCount > 0
                          ? String(historyCurrentCount)
                          : '—';

                      return (
                        <div key={`${item.type}-${item.updated_at}-${index}`} className="rounded-md border bg-muted/20 p-3">
                          <div className="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
                            <div className="flex items-center gap-2">
                              <span className="font-medium capitalize">{item.type}</span>
                              <span
                                className={`inline-flex rounded-full px-2.5 py-1 text-xs font-medium ${getTaskStatusClasses(item.status)}`}
                              >
                                {item.status}
                              </span>
                            </div>
                            <div className="text-xs text-muted-foreground">
                              Updated {formatTaskDate(item.updated_at)}
                            </div>
                          </div>
                          <div className="mt-2 grid gap-2 text-sm md:grid-cols-2 xl:grid-cols-4">
                            <div>
                              <p className="text-xs text-muted-foreground">Started</p>
                              <p>{formatTaskDate(item.started_at)}</p>
                            </div>
                            <div>
                              <p className="text-xs text-muted-foreground">Finished</p>
                              <p>{formatTaskDate(item.finished_at)}</p>
                            </div>
                            <div>
                              <p className="text-xs text-muted-foreground">Progress</p>
                              <p>{historyProgress}</p>
                            </div>
                            <div>
                              <p className="text-xs text-muted-foreground">Next Run</p>
                              <p>{formatTaskDate(item.progress?.next_run_at)}</p>
                            </div>
                          </div>
                          {item.last_error && (
                            <div className="mt-2 rounded-md border border-red-200 bg-red-50 p-2 text-xs text-red-700">
                              {item.last_error}
                            </div>
                          )}
                        </div>
                      );
                    })}
                  </div>
                ) : (
                  <div className="text-sm text-muted-foreground">No recent tasks.</div>
                )}
              </CardContent>
            </Card>
          </div>
        </>
      )}
    </div>
  );
}
