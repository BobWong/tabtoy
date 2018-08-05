using System;
using System.Collections.Generic;

using tabtoy;

public enum enumSwitch
{
    switchClose = 0, // 关闭
    switchOpen = 1, // 开启
}

public partial class language_setDefine
{
    /// <summary> 
    /// 编号
    /// </summary>
    public int ID = 0; 
    
    /// <summary> 
    /// 语言
    /// </summary>
    public string language = ""; 
    
    /// <summary> 
    /// 语言包ID
    /// </summary>
    public int languageID = 0; 
    
    /// <summary> 
    /// 手机默认语言编号
    /// </summary>
    public List<int> telephone = new List<int>(); 
    
    /// <summary> 
    /// 语言开关
    /// </summary>
    public enumSwitch enumSwitch = enumSwitch.switchClose; 


    public static List<language_setDefine> readDatas(DataReader reader)
    {
        int len = reader.ReadInt32();

        int[] id = reader.ReadInt32Array(len);
        string[] language = reader.ReadStringArray(len);
        int[] languageID = reader.ReadInt32Array(len);
        int[][] telephone = reader.ReadInt32Array2(len);
        int[] enumSwitch = reader.ReadInt32Array(len);

        List<language_setDefine> list = new List<language_setDefine>();

        language_setDefine temp;
        for (int i = 0; i < len; i++)
        {
            temp = new language_setDefine();
            list.Add(temp);
            temp.ID = id[i];
            temp.language = language[i];
            temp.languageID = languageID[i];
            temp.telephone = new List<int>(telephone[i]);
            temp.enumSwitch = (enumSwitch)enumSwitch[i];
        }
        return list;
    }
} 

public partial class MapNpc
{
    /// <summary> 
    /// 唯一ID
    /// </summary>
    public int ID = 0; 
    
    /// <summary> 
    /// 地图ID
    /// </summary>
    public int mapID = 0; 

    /// <summary> 
    /// NPC位置
    /// </summary>
    public Vec3 NPCPos = new Vec3(); 
    
    /// <summary> 
    /// 朝向
    /// </summary>
    public float posTo = 0f; 
    
    /// <summary> 
    /// 模型
    /// </summary>
    public string modelID = ""; 

    public static List<MapNpc> readDatas(DataReader reader)
    {
        int len = reader.ReadInt32();
        int[] ID = reader.ReadInt32Array(len);
        int[] mapID = reader.ReadInt32Array(len);
        Vec3[] NPCPos = Vec3.createArray(reader, len);
        float[] posTo = reader.ReadFloatArray(len);
        string[] modelID = reader.ReadStringArray(len);

        List<MapNpc> list = new List<MapNpc>();

        MapNpc temp;
        for (int i = 0; i < len; i++)
        {
            temp = new MapNpc();
            list.Add(temp);
            temp.ID = ID[i];
            temp.mapID = mapID[i];
            temp.NPCPos = NPCPos[i];
            temp.posTo = posTo[i];
            temp.modelID = modelID[i];
        }
        return list;
    }
} 


public partial class Vec3
{
    public float X = 0f; 
    public float Y = 0f; 
    public float Z = 0f; 


    public static Vec3 create(tabtoy.DataReader reader)
    {
        Vec3 ins = new Vec3();
        int len = reader.ReadInt32();
        while (len > 0)
        {
            int tag = reader.ReadInt32();
            len--;
            switch (tag)
            {
                case 0x50000:
                    {
                        ins.X = reader.ReadFloat();
                    }
                    break;
                case 0x50001:
                    {
                        ins.Y = reader.ReadFloat();
                    }
                    break;
                case 0x50002:
                    {
                        ins.Z = reader.ReadFloat();
                    }
                    break;
            }
        }
        return ins;
    }

    public static Vec3[] createArray(tabtoy.DataReader reader,Int32 len){
        Vec3[] temp = new Vec3[len];
        for (int i = 0; i < len; i++)
        {
            temp[i] = create(reader);
        }
        return temp;
    }
} 
