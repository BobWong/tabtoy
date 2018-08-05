using System;
using System.IO;
using System.Text;

namespace tabtoy
{
    enum FieldType
    {
        None    = 0,
	    Int32   = 1,
	    Int64   = 2,
	    UInt32  = 3,
	    UInt64  = 4,
	    Float   = 5,
	    String  = 6,
	    Bool    = 7,
	    Enum    = 8,
	    Struct  = 9,	    
    }

    public delegate void DeserializeHandler<T>(T ins, DataReader reader);

    public class DataReader
    {

        static int sizeUInt32 = sizeof(UInt32);
        static int sizeUInt64 = sizeof(UInt64);
        static int sizeInt32 = sizeof(Int32);
        static int sizeInt64 = sizeof(Int64);
        static int sizeFloat = sizeof(float);
        static int sizeBool = sizeof(bool);
        static int sizeByte = sizeof(Byte);

        BinaryReader _reader;
        long _boundPos  = -1;

        public DataReader(Stream stream )
        {
            _reader = new BinaryReader(stream);
            _boundPos = stream.Length;
        }

        public DataReader(Stream stream, long boundpos)
        {
            _reader = new BinaryReader(stream );
            _boundPos = boundpos;
        }

        public DataReader(DataReader reader, long boundpos )
        {
            _reader = reader._reader;
            _boundPos = boundpos;
        }

        void ConsumeData(int size)
        {          
            if ( !IsDataEnough( size ) )
            {
                throw new Exception(string.Format("Out of struct bound pos:{0}/{1} need size:{2}", _reader.BaseStream.Position,_boundPos, size));
            }            
        }

        bool IsDataEnough(int size)
        {            
            return _reader.BaseStream.Position + size <= _boundPos;
        }

        const int CombineFileVersion = 2;

        public bool ReadHeader( )
        {            
            var tag = ReadString();
            if (tag != "TON")
            {
                return false;
            }

            var ver = ReadInt32();
            if (ver != CombineFileVersion)
            {
                return false;
            }
           
            return true;
        }

        static readonly UTF8Encoding encoding = new UTF8Encoding();

        public int ReadTag()
        {
            if ( IsDataEnough(sizeInt32 ) )
            {
                return ReadInt32( );
            }
            return -1;
        }
   
        public int ReadInt32( )
        {
            ConsumeData(sizeInt32);

            return _reader.ReadInt32();
        }

        public int[] ReadInt32Array(int length){
            ConsumeData(sizeInt32 * length);
            int[] arr = new int[length];
            for (int i = 0; i < length; i++)
            {
                arr[i] = _reader.ReadInt32();
            }
            return arr;
        }

        public int[][] ReadInt32Array2(int length)
        {
            int[][] data = new int[length][];
            for (int i = 0; i < length; i++)
            {
                int len = ReadInt32();
                int[] tmp = ReadInt32Array(len);
                data[i] = tmp;
            }
            return data;
        }

        public long ReadInt64( )
        {
            ConsumeData(sizeInt64);

            return _reader.ReadInt64();
        }


        public long[] ReadInt64Array(int length){
            ConsumeData(sizeInt64 * length);
            long[] arr = new long[length];
            for (int i = 0; i < length; i++)
            {
                arr[i] = _reader.ReadInt64();
            }
            return arr;
        }

        public uint ReadUInt32( )
        {
            ConsumeData(sizeUInt32);

            return _reader.ReadUInt32();
        }

        public uint[] ReadUInt32Array(int length){
            ConsumeData(sizeUInt32 * length);
            uint[] arr = new uint[length];
            for (int i = 0; i < length; i++)
            {
                arr[i] = _reader.ReadUInt32();
            }
            return arr;
        }


        public ulong ReadUInt64( )
        {
            ConsumeData(sizeUInt64);

            return _reader.ReadUInt64();
        }

        public ulong[] ReadUInt64Array(int length){
            ConsumeData(sizeUInt64 * length);
            ulong[] arr = new ulong[length];
            for (int i = 0; i < length; i++)
            {
                arr[i] = _reader.ReadUInt64();
            }
            return arr;
        }


        public float ReadFloat( )
        {
            ConsumeData(sizeFloat);

            return _reader.ReadSingle();
        }


        public float[] ReadFloatArray(int length){
            ConsumeData(sizeFloat * length);
            float[] arr = new float[length];
            for (int i = 0; i < length; i++)
            {
                arr[i] = _reader.ReadSingle();
            }
            return arr;
        }

        public bool ReadBool( )
        {
            ConsumeData(sizeBool);

            return _reader.ReadBoolean();
        }

        public bool[] ReadBoolArray(int length){
            ConsumeData(sizeBool * length);
            bool[] arr = new bool[length];
            for (int i = 0; i < length; i++)
            {
                arr[i] = _reader.ReadBoolean();
            }
            return arr;
        }

        public string ReadString()
        {
            var len = ReadInt32();
            ConsumeData(sizeByte * len);
            return encoding.GetString(_reader.ReadBytes(len));
        }

        /// <summary>
        /// 读取一维字符串数组
        /// </summary>
        /// <returns>The string array.</returns>
        /// <param name="length">Length.</param>
        public string[] ReadStringArray(int length){
            string[] arr = new string[length];
            var len = 0;
            for (int i = 0; i < length; i++)
            {   
                len = ReadInt32();
                arr[i] = encoding.GetString(_reader.ReadBytes(len));
            }
            return arr;
        }

        /// <summary>
        /// 读取二维字符串数组
        /// </summary>
        /// <returns>The string array2.</returns>
        /// <param name="length">Length.</param>
        public string[][] ReadStringArray2(int length)
        {
            string[][] data = new string[length][];
            for (int i = 0; i < length; i++)
            {
                int len = ReadInt32();
                string[] tmp = ReadStringArray(len);
                data[i] = tmp;
            }
            return data;
        }

        public T ReadEnum<T>( )
        {
            return (T)Enum.ToObject(typeof(T), ReadInt32());                
        }

        public T ReadStruct<T>(DeserializeHandler<T> handler) where T : class
        {
            var bound = _reader.ReadInt32();

            var element = Activator.CreateInstance<T>();

            handler(element, new DataReader(this, _reader.BaseStream.Position + bound));

            return element;
        }
        
    }
}
