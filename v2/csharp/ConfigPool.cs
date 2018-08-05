/////////////////////////////////////////////////////////////////////////////////
// @desc 静态表读取池
// @copyright ©2017
// @release 2018年1月11日 星期四
// @author Byhy
/////////////////////////////////////////////////////////////////////////////////

using System;
using System.IO;
using System.Collections.Generic;
using tabtoy;
using Pis;

public class ConfigPool {

    private static Dictionary<string, IConfig> m_s_cache = new Dictionary<string, IConfig>();

    /// <summary>
    /// 清理缓存
    /// </summary>
    public static void ClearCache() {
        m_s_cache.Clear();
    }

    // 获取
    public static T Get<T>(string path) where T : IConfig {
        Type byteClassType = typeof(T);
        string name = byteClassType.ToString().Substring(6);

        if (m_s_cache.ContainsKey(name)) {
            return (T)m_s_cache[name];
        }

        byte[] configBytes = File.ReadAllBytes(string.Format("{0}/{1}.bytes",path ,name));

        if (configBytes.Length > 0)
        {
            using (MemoryStream stream = new MemoryStream(configBytes))
            {
                stream.Position = 0;
                var reader = new DataReader(stream, stream.Length);
                if (!reader.ReadHeader())
                {
                    Console.WriteLine(">>>>tabtoy: {0}, combine file crack!");
                }
                else
                {
                    Console.WriteLine(name + " data len :" + configBytes.Length);
                }
                T result;
                var data = Activator.CreateInstance(byteClassType);

                var mInfos = byteClassType.GetMethods();
                foreach (System.Reflection.MethodInfo info in mInfos)
                {
                    if (info.Name == "Deserialize")
                    {
                        result = (T)info.Invoke(data, new object[] { reader });
                        m_s_cache.Add(name, result);
                        return result;
                    }
                }
            }
        }
        return default(T);

    }

}